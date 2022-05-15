package alerts

import (
	"strconv"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

var Alert *DiscordHook

// DiscordHook : Struct to interact with discord hook
type DiscordHook struct {
	Disabled bool
	Title    string

	client webhook.Client

	iD    string
	token string
}

// SendEmbed : Self Explainatory
func (g *DiscordHook) SendEmbed(dat string, feilds map[string]string) error {

	var embed *discord.EmbedBuilder = discord.NewEmbedBuilder()

	if len(feilds) > 0 {
		for k, v := range feilds {
			embed.AddField(k, v, false)
		}
	}

	arr := FormatMsg(dat)

	if len(arr) < 5 {
		var er error = nil
		for _, v := range arr {

			e := discord.Embed{
				Title:       g.Title,
				Description: v,
				Type:        discord.EmbedTypeArticle,

				Fields: embed.Fields,
				Color:  NormColor,
			}

			_, er = g.client.CreateEmbeds([]discord.Embed{e})

		}
		return er
	} else {
		wg := &sync.WaitGroup{}
		m := &sync.Mutex{}
		var er error

		for _, v := range arr {
			wg.Add(1)
			go func(x string) {
				defer wg.Done()
				e := discord.Embed{
					Title:       g.Title,
					Description: x,
					Type:        discord.EmbedTypeArticle,

					Fields: embed.Fields,
					Color:  NormColor,
				}

				_, z := g.client.CreateEmbeds([]discord.Embed{e})
				if z != nil {
					m.Lock()
					er = z
					m.Unlock()
				}
			}(v)
		}

		wg.Wait()

		return er

	}

}

// SendErr : Self Explainatory
func (g *DiscordHook) SendErr(er error, cmd string) error {

	f := discord.NewEmbedBuilder()
	f.AddField("CMD", cmd, false)

	e := discord.Embed{
		Title:       g.Title,
		Description: er.Error(),
		Type:        discord.EmbedTypeArticle,

		Fields: f.Fields,
		Color:  ErrColor,
	}

	_, er = g.client.CreateEmbeds([]discord.Embed{e})

	return er

}

// NewDiscordHook :
func NewDiscordHook(id string, token string) *DiscordHook {

	webid := snowflake.ID(0)

	no, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}

	webid = snowflake.ID(no)

	clientx := webhook.NewClient(webid, token)

	g := DiscordHook{
		client: clientx,
		token:  token,
		iD:     id,
	}

	rest.WithQueryParam("splitlines", "true")

	return &g
}
