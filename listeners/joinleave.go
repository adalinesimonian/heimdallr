package listeners

import (
	"github.com/cbroglie/mustache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/myrkvi/heimdallr/model"
	"github.com/myrkvi/heimdallr/utils"
	"log/slog"
)

func OnUserJoin(e *events.GuildMemberJoin) {
	guildID := e.GuildID
	guild, err := e.Client().Rest().GetGuild(guildID, false)
	if err != nil {
		return
	}

	guildSettings, err := model.GetGuildSettings(guildID)
	if err != nil {
		return
	}

	if !guildSettings.JoinMessageEnabled {
		return
	}

	joinLeaveChannel := guildSettings.JoinLeaveChannel
	if joinLeaveChannel == 0 && guild.SystemChannelID != nil {
		joinLeaveChannel = *guild.SystemChannelID
	}
	if joinLeaveChannel == 0 {
		return
	}

	joinleaveInfo := utils.NewMessageTemplateData(e.Member, guild.Guild)

	contents, err := mustache.RenderRaw(guildSettings.JoinMessage, true, joinleaveInfo)
	if err != nil {
		slog.Error("Failed to render join message template.",
			"err", err,
			"guild_id", guildID,
		)
		return
	}

	_, err = e.Client().Rest().CreateMessage(joinLeaveChannel, discord.NewMessageCreateBuilder().SetContent(contents).Build())
}

func OnUserLeave(e *events.GuildMemberLeave) {
	guildID := e.GuildID
	guild, err := e.Client().Rest().GetGuild(guildID, false)
	if err != nil {
		return
	}

	guildSettings, err := model.GetGuildSettings(guildID)
	if err != nil {
		return
	}

	if !guildSettings.LeaveMessageEnabled {
		return
	}

	joinLeaveChannel := guildSettings.JoinLeaveChannel
	if joinLeaveChannel == 0 && guild.SystemChannelID != nil {
		joinLeaveChannel = *guild.SystemChannelID
	}
	if joinLeaveChannel == 0 {
		return
	}

	joinleaveInfo := utils.NewMessageTemplateData(e.Member, guild.Guild)

	contents, err := mustache.RenderRaw(guildSettings.LeaveMessage, true, joinleaveInfo)
	if err != nil {
		slog.Error("Failed to render leave message template.",
			"err", err,
			"guild_id", guildID,
		)
		return
	}

	_, err = e.Client().Rest().CreateMessage(joinLeaveChannel, discord.NewMessageCreateBuilder().SetContent(contents).Build())
}
