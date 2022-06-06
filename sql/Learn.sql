USE [kequ5060]
GO

/****** Object:  Table [dbo].[zbot_learn]    Script Date: 2022/6/6 10:56:45 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[zbot_learn](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[ask] [varchar](255) NOT NULL,
	[group_id] [bigint] NOT NULL,
	[admin_id] [bigint] NULL,
	[answer] [varchar](max) NULL,
	[gmt_modified] [datetime2](7) NULL,
	[pass] [bit] NULL
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO

