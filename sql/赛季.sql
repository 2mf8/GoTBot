USE [kequ5060]
GO

/****** Object:  Table [dbo].[guild_achievement]    Script Date: 2023/8/12 19:38:10 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[guild_achievement](
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[user_id] [varchar](32) NOT NULL,
	[user_name] [nvarchar](32) NOT NULL,
	[avatar] [varchar](128) NULL,
	[item] [varchar](16) NOT NULL,
	[best] [int] NOT NULL,
	[average] [int] NOT NULL,
	[session] [int] NOT NULL,
 CONSTRAINT [PK_guild_achievement] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

