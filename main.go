package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aikchun/aqicnsg-telegram/aqicnsg"
	"github.com/aikchun/gotelegrambot"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// SendMessageParams holds the parameters of the Telegram Bot API sendMessage
// method (https://core.telegram.org/bots/api#sendmessage). Nested objects are
// left as raw JSON.
type SendMessageParams struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	// ChatID is either an integer chat id or a "@username" string.
	ChatID                     any             `json:"chat_id"`
	MessageThreadID            int64           `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID      int64           `json:"direct_messages_topic_id,omitempty"`
	EphemeralMessageParameters json.RawMessage `json:"ephemeral_message_parameters,omitempty"`
	Text                       string          `json:"text"`
	ParseMode                  string          `json:"parse_mode,omitempty"`
	Entities                   json.RawMessage `json:"entities,omitempty"`
	LinkPreviewOptions         json.RawMessage `json:"link_preview_options,omitempty"`
	DisableNotification        bool            `json:"disable_notification,omitempty"`
	ProtectContent             bool            `json:"protect_content,omitempty"`
	AllowPaidBroadcast         bool            `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID            string          `json:"message_effect_id,omitempty"`
	SuggestedPostParameters    json.RawMessage `json:"suggested_post_parameters,omitempty"`
	ReplyParameters            json.RawMessage `json:"reply_parameters,omitempty"`
	ReplyMarkup                json.RawMessage `json:"reply_markup,omitempty"`
}

type botParams struct {
	botToken string
	params   map[string]string
}

// commandHandler needs to have an Update and the botToken
type CommandHandler = func(u gotelegrambot.Update, bp botParams)

type bot struct {
	botToken        string
	commandsMap     map[string]CommandHandler
	fallbackHandler CommandHandler
	params          map[string]string
}

func (b *bot) register(commandName string, handler CommandHandler) {
	if b.commandsMap == nil {
		b.commandsMap = map[string]CommandHandler{}
	}
	b.commandsMap[commandName] = handler
}

func (b *bot) registerFallbackHandler(handler CommandHandler) {
	b.fallbackHandler = handler
}

func (b *bot) getCommandHandler(u *gotelegrambot.Update) CommandHandler {
	command, _, ok := parseText(u.Message.Text)
	if !ok {
		log.Printf("unable to parse send message text: %s", u.Message.Text)
		return b.fallbackHandler
	}

	commandHandler, ok := b.commandsMap[command]
	if !ok {
		if b.fallbackHandler != nil {
			log.Printf("unable to parse command from text: %s", u.Message.Text)
			return b.fallbackHandler
		}
	}
	return commandHandler

}

func (b *bot) handleUpdate(u *gotelegrambot.Update) {
	commandHandler := b.getCommandHandler(u)
	commandHandler(*u, botParams{botToken: b.botToken, params: maps.Clone(b.params)})
}

// aqiEmoji returns the coloured circle for the AQI band the value falls in.
func aqiEmoji(aqi int) string {
	switch {
	case aqi <= 50:
		return "🟢"
	case aqi <= 100:
		return "🟡"
	case aqi <= 150:
		return "🟠"
	case aqi <= 200:
		return "🔴"
	case aqi <= 300:
		return "🟣"
	default:
		return "⚫"
	}
}

// parseCommand extracts the command from message text of the form
// "/{command}@{botname}". The "@{botname}" part and anything after the first
// space are optional and dropped. It returns "" if text is not a command.
func parseText(text string) (string, []string, bool) {
	tokens := strings.Split(text, " ")
	text, ok := strings.CutPrefix(strings.TrimSpace(tokens[0]), "/")
	if !ok {
		return "", []string{}, ok
	}
	text, _, _ = strings.Cut(text, " ")
	command, _, _ := strings.Cut(text, "@")
	return command, tokens[1:], true
}

// buildPM25Report fetches the current area AQIs and formats them as the
// PM2.5 report text.
func buildPM25Report(apiToken string) (string, error) {
	areas, err := getAreaAqis(apiToken)
	if err != nil {
		return "", err
	}
	return formatPM25Report(areas)
}

// getAreaAqis fetches the current area AQIs, failing if none are returned.
func getAreaAqis(apiToken string) ([]aqicnsg.AreaAqi, error) {
	areas, err := aqicnsg.GetAreaAqis(apiToken)
	if err != nil {
		return nil, fmt.Errorf("get area aqis: %w", err)
	}
	if len(areas) == 0 {
		return nil, fmt.Errorf("get area aqis: no areas returned")
	}
	return areas, nil
}

// maxAQI returns the highest AQI across areas.
func maxAQI(areas []aqicnsg.AreaAqi) int {
	m := areas[0].AQI
	for _, a := range areas[1:] {
		m = max(m, a.AQI)
	}
	return m
}

// formatPM25Report formats areas as the PM2.5 report text.
func formatPM25Report(areas []aqicnsg.AreaAqi) (string, error) {
	lines := make([]string, len(areas))
	for i, a := range areas {
		city, _, _ := strings.Cut(a.City, ",")

		lines[i] = fmt.Sprintf("%s: %d %s", city, a.AQI, aqiEmoji(a.AQI))
	}
	timeObj, err := time.Parse(time.DateTime, areas[0].Time)
	if err != nil {
		return "", fmt.Errorf("time format err: %w", err)
	}
	return fmt.Sprintf("PM2.5 AQI at %s %d %s:\n%s", timeObj.Month().String()[:3], timeObj.Day(), timeObj.Format(time.Kitchen), strings.Join(lines, "\n")), nil
}

// sendMessage sends text to chatID via the Telegram Bot API.
func sendMessage(botToken string, chatID int64, text string) error {
	b, err := json.Marshal(SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("marshal send message params: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("send message: %s: %s", resp.Status, body)
	}
	return nil
}

func handleNowCommand(u gotelegrambot.Update, bp botParams) {
	text, err := buildPM25Report(bp.params["AQICN_API_TOKEN"])
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if err := sendMessage(bp.botToken, u.Message.Chat.ID, text); err != nil {
		log.Printf("%v", err)
	}
}

func handleSubscribe(u gotelegrambot.Update, bp botParams) {
	conn, err := pgx.Connect(context.Background(), bp.params["DATABASE_URL"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	chatID := u.Message.Chat.ID
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	var ID int64
	err = conn.QueryRow(context.Background(), "select id from subscription where chat_id=$1", chatID).Scan(&ID)
	if err != nil {
		insert := fmt.Sprintf("INSERT INTO subscription(chat_id, created_at, updated_at) values (%d, '%s', '%s')", chatID, now, now)
		err = pgx.BeginFunc(context.Background(), conn, func(tx pgx.Tx) error {
			_, err := tx.Exec(context.Background(), insert)
			return err
		})
		if err != nil {
			log.Printf("unable to insert subscription %v", err)
			log.Printf("insert: %s", insert)
			return
		}
		return
	}
	err = pgx.BeginFunc(context.Background(), conn, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), fmt.Sprintf("UPDATE subscription SET is_active = %t, updated_at = '%s' WHERE id = %d", true, now, ID))
		return err
	})
	if err != nil {
		log.Printf("unable to update subscription %v", err)
		return
	}

}

func handleUnsubscribe(u gotelegrambot.Update, bp botParams) {
	conn, err := pgx.Connect(context.Background(), bp.params["DATABASE_URL"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	chatID := u.Message.Chat.ID
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	err = pgx.BeginFunc(context.Background(), conn, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), fmt.Sprintf("UPDATE subscription SET is_active = %t, updated_at = '%s' WHERE chat_id = %d", false, now, chatID))
		return err
	})
	if err != nil {
		log.Printf("unable to update subscription %v", err)
		return
	}

}

func handleSendReportLambda(u gotelegrambot.Update) {
	bp := botParams{
		botToken: os.Getenv("BOT_TOKEN"),
		params: map[string]string{
			"DATABASE_URL":    os.Getenv("DATABASE_URL"),
			"AQICN_API_TOKEN": os.Getenv("AQICN_API_TOKEN"),
		},
	}

	handleSendReport(u, bp)
}
func handleSendReport(_ gotelegrambot.Update, bp botParams) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, bp.params["DATABASE_URL"])
	if err != nil {
		log.Printf("unable to connect to database: %v", err)
		return
	}
	defer conn.Close(ctx)

	// Fetch the AQIs once and send the same text to every subscriber.
	areas, err := getAreaAqis(bp.params["AQICN_API_TOKEN"])
	if err != nil {
		log.Printf("%v", err)
		return
	}

	// Only notify subscribers whose threshold is exceeded by at least one area.
	rows, err := conn.Query(ctx, "SELECT chat_id FROM subscription WHERE is_active = true AND aqi_threshold < $1", maxAQI(areas))
	if err != nil {
		log.Printf("unable to query subscriptions: %v", err)
		return
	}
	chatIDs, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		log.Printf("unable to read subscriptions: %v", err)
		return
	}
	if len(chatIDs) == 0 {
		return
	}

	text, err := formatPM25Report(areas)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	for _, chatID := range chatIDs {
		if err := sendMessage(bp.botToken, chatID, text); err != nil {
			log.Printf("chat %d: %v", chatID, err)
		}
	}
}

func handleCommandNotFound(u gotelegrambot.Update, bp botParams) {
	b, err := json.Marshal(SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "No such function",
	})
	if err != nil {
		log.Printf("marshal send message params: %v", err)
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bp.botToken)

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Printf("send message: %v", err)
		return
	}
	defer resp.Body.Close()

}

func handleUpdate(u gotelegrambot.Update) {
	// determine what command the Update wants to call
	b := bot{botToken: os.Getenv("BOT_TOKEN")}
	b.params = map[string]string{}
	b.params["DATABASE_URL"] = os.Getenv("DATABASE_URL")
	b.params["AQICN_API_TOKEN"] = os.Getenv("AQICN_API_TOKEN")
	b.register("now", handleNowCommand)
	b.register("subscribe", handleSubscribe)
	b.register("unsubscribe", handleUnsubscribe)
	b.registerFallbackHandler(handleCommandNotFound)
	b.handleUpdate(&u)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var u gotelegrambot.Update

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &u); err != nil {
		log.Fatal(err)
	}

	handleUpdate(u)

}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Couldn't find .env")
	}

	e := os.Getenv("ENVIRONMENT")

	if e == "dev" {
		p := ":8080"
		http.HandleFunc("/", handler)
		fmt.Printf("starting http server\n")
		fmt.Printf("listening: %s\n", p)

		log.Fatal(http.ListenAndServe(p, nil))
	}

	if e == "cron" {
		lambda.Start(handleSendReportLambda)
	}
	// handleSendReportLambda(gotelegrambot.Update{})

	if e == "prod" {
		lambda.Start(handleUpdate)
	}

}
