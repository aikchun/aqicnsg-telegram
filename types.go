package main

import (
	"encoding/json"

	"github.com/aikchun/gotelegrambot"
)

type botParams struct {
	botToken string
	params   map[string]string
}

// commandHandler needs to have an Update and the botToken
type CommandHandler = func(u gotelegrambot.Update, bp botParams) error

type bot struct {
	botToken        string
	commandsMap     map[string]CommandHandler
	fallbackHandler CommandHandler
	params          map[string]string
}

// SendMessageParams holds the parameters of the Telegram Bot API sendMessage
// method (https://core.telegram.org/bots/api#sendmessage). Nested objects other
// than ReplyMarkup are left as raw JSON.
type SendMessageParams struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	// ChatID is either an integer chat id or a "@username" string.
	ChatID                     any                   `json:"chat_id"`
	MessageThreadID            int64                 `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID      int64                 `json:"direct_messages_topic_id,omitempty"`
	EphemeralMessageParameters json.RawMessage       `json:"ephemeral_message_parameters,omitempty"`
	Text                       string                `json:"text"`
	ParseMode                  string                `json:"parse_mode,omitempty"`
	Entities                   json.RawMessage       `json:"entities,omitempty"`
	LinkPreviewOptions         json.RawMessage       `json:"link_preview_options,omitempty"`
	DisableNotification        bool                  `json:"disable_notification,omitempty"`
	ProtectContent             bool                  `json:"protect_content,omitempty"`
	AllowPaidBroadcast         bool                  `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID            string                `json:"message_effect_id,omitempty"`
	SuggestedPostParameters    json.RawMessage       `json:"suggested_post_parameters,omitempty"`
	ReplyParameters            json.RawMessage       `json:"reply_parameters,omitempty"`
	ReplyMarkup                *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// InlineKeyboardMarkup is an inline keyboard that appears right next to the
// message it belongs to
// (https://core.telegram.org/bots/api#inlinekeyboardmarkup).
type InlineKeyboardMarkup struct {
	// InlineKeyboard is an array of button rows.
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	ForceReply     bool                     `json:"force_reply,omitempty"`
}

// InlineKeyboardButton is one button of an inline keyboard
// (https://core.telegram.org/bots/api#inlinekeyboardbutton). Exactly one of
// the fields other than Text, IconCustomEmojiID and Style must be set.
type InlineKeyboardButton struct {
	Text              string `json:"text"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
	// Style is one of "danger", "success" or "primary".
	Style        string      `json:"style,omitempty"`
	URL          string      `json:"url,omitempty"`
	CallbackData string      `json:"callback_data,omitempty"`
	WebApp       *WebAppInfo `json:"web_app,omitempty"`
	LoginURL     *LoginURL   `json:"login_url,omitempty"`
	// SwitchInlineQuery and SwitchInlineQueryCurrentChat may be set to an
	// empty string, so they are pointers to distinguish that from unset.
	SwitchInlineQuery            *string                      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string                      `json:"switch_inline_query_current_chat,omitempty"`
	SwitchInlineQueryChosenChat  *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
	CopyText                     *CopyTextButton              `json:"copy_text,omitempty"`
	CallbackGame                 *CallbackGame                `json:"callback_game,omitempty"`
	Pay                          bool                         `json:"pay,omitempty"`
	Disabled                     *DisabledButton              `json:"disabled,omitempty"`
}

// WebAppInfo describes a Web App
// (https://core.telegram.org/bots/api#webappinfo).
type WebAppInfo struct {
	URL string `json:"url"`
}

// LoginURL is used by an inline keyboard button to automatically authorize a
// user (https://core.telegram.org/bots/api#loginurl).
type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text,omitempty"`
	BotUsername        string `json:"bot_username,omitempty"`
	RequestWriteAccess bool   `json:"request_write_access,omitempty"`
}

// SwitchInlineQueryChosenChat switches the user to inline mode in a chosen
// chat (https://core.telegram.org/bots/api#switchinlinequerychosenchat).
type SwitchInlineQueryChosenChat struct {
	Query             string `json:"query,omitempty"`
	AllowUserChats    bool   `json:"allow_user_chats,omitempty"`
	AllowBotChats     bool   `json:"allow_bot_chats,omitempty"`
	AllowGroupChats   bool   `json:"allow_group_chats,omitempty"`
	AllowChannelChats bool   `json:"allow_channel_chats,omitempty"`
}

// CopyTextButton copies Text to the clipboard
// (https://core.telegram.org/bots/api#copytextbutton).
type CopyTextButton struct {
	Text string `json:"text"`
}

// CallbackGame is a placeholder that currently holds no information
// (https://core.telegram.org/bots/api#callbackgame).
type CallbackGame struct{}

// DisabledButton marks a button as disabled; it currently holds no
// information (https://core.telegram.org/bots/api#disabledbutton).
type DisabledButton struct{}
