package main

import (
	"encoding/json"
	"fmt"
	"github.com/NicoNex/echotron"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type bot struct {
	chatId int64
	echotron.Api
	roundRiri int
}

const (
	botLogsFolder  = "/NinoBotData/logs/"
	botAudioFolder = "/NinoBotData/audio/"
)

var TOKEN = os.Getenv("NinoBot")
var logsFolder string
var audioFolder string

func newBot(chatId int64) echotron.Bot {
	return &bot{
		chatId,
		echotron.NewApi(TOKEN),
		0,
	}
}

func (b *bot) makeButtons(buttonsText []string, callbacksData []string, layout int) ([]byte, error) {
	if layout != 1 && layout != 2 {
		return nil, fmt.Errorf("wrong layout")
	}
	if len(buttonsText) != len(callbacksData) {
		return nil, fmt.Errorf("different text and data length")
	}

	buttons := make([]echotron.InlineButton, 0)
	for i, v := range buttonsText {
		buttons = append(buttons, b.InlineKbdBtn(v, "", callbacksData[i]))
	}

	keys := make([]echotron.InlineKbdRow, 0)
	switch layout {
	case 1:
		for i := 0; i < len(buttons); i++ {
			keys = append(keys, echotron.InlineKbdRow{buttons[i]})
		}
		break
	case 2:
		for i := 0; i < len(buttons); i += 2 {
			if i+1 < len(buttons) {
				keys = append(keys, echotron.InlineKbdRow{buttons[i], buttons[i+1]})
			} else {
				keys = append(keys, echotron.InlineKbdRow{buttons[i]})
			}
		}
		break
	}

	inlineKMarkup := b.InlineKbdMarkup(keys...)
	return inlineKMarkup, nil
}

func initFolders() {
	currentPath, _ := os.Getwd()

	logsFolder = currentPath + botLogsFolder
	_ = os.MkdirAll(logsFolder, 0755)

	audioFolder = currentPath + botAudioFolder
	_ = os.MkdirAll(audioFolder, 0755)
}

func main() {
	initFolders()

	dsp := echotron.NewDispatcher(TOKEN, newBot)
	dsp.ListenWebhook("https://hiddenfile.tk:443/bot/NinoBot", 40991)
}

func (b *bot) Update(update *echotron.Update) {
	b.logUser(update, logsFolder)
	if update.Message != nil {
		messageTextLower := strings.ToLower(update.Message.Text)
		if messageTextLower == "/start" {
			b.sendStart(update.Message)
		} else if messageTextLower == "/credits" {
			b.sendCredits(update)
		} else if b.roundRiri == 1 {
			b.sendNinoTypicalExpression(update.Message, 7)
		} else if strings.Contains(messageTextLower, "ball") && (strings.Contains(messageTextLower, "nino") ||
			strings.Contains(messageTextLower, "ni") || strings.Contains(messageTextLower, "nì")) {
			b.sendNinoTypicalExpression(update.Message, 30)
		} else if strings.Contains(messageTextLower, "cant") && strings.Contains(messageTextLower, "canzone") {
			b.sendNinoSong(update.Message)
		} else if strings.Contains(messageTextLower, "anni") && (strings.Contains(messageTextLower, "nino") ||
			strings.Contains(messageTextLower, "ni") || strings.Contains(messageTextLower, "nì")) {
			b.sendNinoTypicalExpression(update.Message, 1)
		} else if strings.Contains(messageTextLower, "ciao") && strings.Contains(messageTextLower, "nino") {
			b.sendNinoTypicalExpression(update.Message, 10)
		} else if strings.Contains(messageTextLower, "buon") && (strings.Contains(messageTextLower, "giorn") ||
			strings.Contains(messageTextLower, "sera")) {
			b.sendNinoTypicalExpression(update.Message, 10)
		} else if strings.Contains(messageTextLower, "donna") && strings.Contains(messageTextLower, "ideale") {
			b.sendNinoTypicalExpression(update.Message, -1)
		} else if strings.Contains(messageTextLower, "nino") || strings.Contains(messageTextLower, "nì") {
			b.sendNinoTypicalExpression(update.Message, -1)
		} else if update.Message.Chat.Type == "private" {
			b.privateTalkWithNino(update.Message)
		}
	} else if update.Message == nil && update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "credits" {
			b.sendCredits(update)
		}
	}

}

func (b *bot) sendCredits(update *echotron.Update) {
	var chatId int64
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		chatId = update.Message.Chat.ID
	}

	b.SendMessage("🤖 Bot creato da @GiovanniRanaTortello\n😺 GitHub: https://github.com/DarkFighterLuke\n"+
		"\n🌐 Proudly hosted on Raspberry Pi 3\n"+
		"\nContribuisci anche tu alla linguistica di NinoBot su GitHub o contattando il creatore!\n"+
		"N.B. Questo bot è satirico e non intende offendere chi rappresenta. "+
		"Ti auguriamo di trovare l'amore Nino.", chatId, echotron.PARSE_HTML)
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(update.CallbackQuery.ID, "Crediti", false)
	}
}

func (b *bot) sendStart(message *echotron.Message) {
	msg := `<b>Hai contattato Nino!</b>
Piacere di conoscerti, %s!
Io sono Nino.
Ho 47 anni. Sono di Paceco.
`
	msg = fmt.Sprintf(msg, message.User.FirstName)

	buttonText := []string{"Credits 🌟"}
	buttonCallback := []string{"credits"}
	buttons, err := b.makeButtons(buttonText, buttonCallback, 1)
	if err != nil {
		log.Println("Error creating buttons:", err)
	}

	b.SendMessageWithKeyboard(msg, message.Chat.ID, buttons, echotron.PARSE_HTML)
}

func (b *bot) sendNinoTypicalExpression(message *echotron.Message, n int) {
	if n < 0 {
		n = rand.Intn(31)
	}

	switch n {
	case 0:
		msg := "Io sono Nino."
		b.SendMessage(msg, message.Chat.ID)
		break
	case 1:
		msg := "Ho 47 anni."
		b.SendMessage(msg, message.Chat.ID)
		break
	case 2:
		msg := "Sono di Paceco."
		b.SendMessage(msg, message.Chat.ID)
		break
	case 3:
		msg := "A te u musu ti fazzu shcattare Ngichinè!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 4:
		msg := "A ballare tanto assai non mi piace"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 5:
		msg := "*risata ebete*"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 6:
		b.roundRiri = 1
		msg := "Ma che cazzo ci riri oh!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 7:
		b.roundRiri = 0
		msg := "Ma che caspitina ci riri!!!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 8:
		msg := "Cosa volete fare facete"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 9:
		msg := "Se volete chiamare chiamate"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 10:
		msg := "Buonasera."
		b.SendMessage(msg, message.Chat.ID)
		break
	case 11:
		msg := "Non ho stato mai fidanzato"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 12:
		msg := "Vuai ceccando un'anima gemella"
		b.SendMessage("...", message.Chat.ID)
		b.SendMessage(msg, message.Chat.ID)
		break
	case 13:
		// Sticker "Buonasera"
		stickerId := "CAACAgQAAxkBAANlYDP628W4thKmIIM2TktXp3n0QOIAAnoAA5tcdge1GHiyda2EVx4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 14:
		// Sticker "Io sono Nino."
		stickerId := "CAACAgQAAxkBAANnYDP68R5djBgo6jDyxMYVcpQ3yi0AAnsAA5tcdgcxQ_AoiEAefR4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 15:
		// Sends "A_ballare_tanto_assai_non_mi_piace.mp3"
		fileId := "AwACAgQAAxkBAANbYDLybzTgltgkHe2e2lpQX6bmwfEAAo8IAALWIZlRfjem1-W3x64eBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 16:
		// Sends "Che_caspitina_ci_ridi.mp3"
		fileId := "AwACAgQAAxkBAANZYDLyVJ-p19HnAAG-Izv7W1HkjfDsAAKOCAAC1iGZUWA1lqOtsg3DHgQ"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 17:
		// Sends "Cosa_volete_fare_facete_se_volete_chiamare_chiamate.mp3"
		fileId := "AwACAgQAAxkBAANhYDLzeDpELsPcpnZxhES7EKi1mQkAApoIAALWIZlREY6FUQS42CMeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 18:
		// Sends "Io_sono_Nino.mp3"
		fileId := "AwACAgQAAxkBAANTYDLx8E7HybhXzoD8evbXLisAASw9AAKLCAAC1iGZUXj422CmkJ5iHgQ"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 19:
		// Sends "Pescare.mp3"
		fileId := "AwACAgQAAxkBAANVYDLyF-hUXgbOLjtfssLdciZbwKQAAowIAALWIZlRfj-n3E5V2BEeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 20:
		// Sends "U_musu_t_fazzu_shcattare.mp3"
		fileId := "AwACAgQAAxkBAANXYDLyNDT7RLcfp6S-bsEWswqShhoAAo0IAALWIZlRtxc1za0ltWoeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 21:
		// Sends "Vado_cercando_un_anima_gemella.mp3"
		fileId := "AwACAgQAAxkBAANdYDLzHHZ82eZ0O6LNxmZGqzmysOwAApYIAALWIZlR_z6E4_-svTQeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 22:
		// Sticker "U musu t fazzu shcattare Ngichinè!"
		stickerId := "CAACAgQAAxkBAANpYDP7Dp1Cmbj6bM9U4sVStWDfgHAAAnwAA5tcdgekswjUH7G0Zh4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 23:
		// Sticker "E riri che cazzu ci riri oh?!"
		stickerId := "CAACAgQAAxkBAANrYDP7Id0Y63nl49N-uDKMNg_ZWnQAAn0AA5tcdgelNbWrYhus4h4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 24:
		// Sticker "Ma che caspitina ci ridi!"
		stickerId := "CAACAgQAAxkBAANtYDP7NN0IowqhqPoChpN5VvukYMQAAn4AA5tcdgf90-GacfV7_B4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 25:
		// Sticker "A ballare tanto assai non mi piace"
		stickerId := "CAACAgQAAxkBAANvYDP7SLA_OQJulmt3DC3_1OxWiNYAAn8AA5tcdgetgSpHh9qWQB4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 26:
		// Sticker "...mangiare una pizza..."
		stickerId := "CAACAgQAAxkBAANxYDP7ZpQbIbUvqFlvAjCrp5TAHRUAAoAAA5tcdgcdPMUP3VtDTx4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 27:
		// Sticker "Che cosa volete fare facete"
		stickerId := "CAACAgQAAxkBAANzYDP7dxchhkHIAAH-Ci5tvi1VkBLiAAKBAAObXHYHeS-dZVQDxz0eBA"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 28:
		// Sticker "Se volete chiamare chiamate"
		stickerId := "CAACAgQAAxkBAAN1YDP7jID-G7wiYFKvrcbaZQm_ZYkAAoIAA5tcdgeP02ixTce2Ix4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 29:
		// Sticker "Nino."
		stickerId := "CAACAgQAAxkBAAN3YDQDXKq2xT1JupwouKkjUo34PZwAAoMAA5tcdgdIQbmeKDiqmx4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 30:
		// Sends "Nino danzante" GIF
		gifId := "CgACAgQAAxkBAAN7YDQjCwKIGY2GfnuIbLiYEru-HTwAAiUNAAK-UqFRgCc7-tmphdoeBA"
		b.SendAnimationByID(gifId, message.Chat.ID)
	}
}

func (b *bot) privateTalkWithNino(message *echotron.Message) {
	n := rand.Float32()
	if n < 0.75 {
		b.sendNinoTypicalExpression(message, -1)
	}
}

func (b *bot) sendNinoSong(message *echotron.Message) {
	n := rand.Intn(1)

	switch n {
	case 0:
		fileId := "AwACAgQAAxkBAANjYDLz7kzO7KtS9lFw0xWs_CAvhvAAApsIAALWIZlR1gl30_hxQ7YeBA"
		b.SendVoiceByID(fileId, "By Davide Belvedere", message.Chat.ID)
		break
	}
}

func (b *bot) logUser(update *echotron.Update, folder string) {
	data, err := json.Marshal(update)
	if err != nil {
		log.Println("Error marshaling logs: ", err)
		return
	}

	var filename string

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message.Chat.Type == "private" {
			if update.CallbackQuery.Message.Chat.Username == "" {
				filename = folder + update.CallbackQuery.Message.Chat.FirstName + "_" + update.CallbackQuery.Message.Chat.LastName + ".txt"
			} else {
				filename = folder + update.CallbackQuery.Message.Chat.Username + ".txt"
			}
		} else {
			filename = folder + update.Message.Chat.Title + ".txt"
		}

	} else if update.Message != nil {
		if update.Message.Chat.Type == "private" {
			if update.Message.Chat.Username == "" {
				filename = folder + update.Message.Chat.FirstName + "_" + update.Message.Chat.LastName + ".txt"
			} else {
				filename = folder + update.Message.Chat.Username + ".txt"
			}
		} else {
			filename = folder + update.Message.Chat.Title + ".txt"
		}

	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	dataString := time.Now().Format("2006-01-02T15:04:05") + string(data[:])
	_, err = f.WriteString(dataString + "\n")
	if err != nil {
		log.Println(err)
		return
	}
	err = f.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
