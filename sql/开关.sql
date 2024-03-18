USE [kequ5060]
GO

/****** Object:  Table [dbo].[guild_switch]    Script Date: 2023/8/12 19:36:44 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[guild_switch](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[guild_id] [varchar](32) NOT NULL,
	[channel_id] [varchar](32) NOT NULL,
	[is_close_or_guard] [bigint] NOT NULL,
	[admin_id] [varchar](32) NULL,
	[gmt_modified] [datetime2](7) NULL,
 CONSTRAINT [PK_guild_switch] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

