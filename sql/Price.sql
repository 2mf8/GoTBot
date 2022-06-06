USE [kequ5060]
GO

/****** Object:  Table [dbo].[zbot_price]    Script Date: 2022/6/6 10:54:26 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[zbot_price](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[group_id] [bigint] NOT NULL,
	[brand] [varchar](100) NULL,
	[item] [varchar](100) NOT NULL,
	[price] [varchar](100) NULL,
	[shipping] [varchar](100) NULL,
	[updater] [bigint] NOT NULL,
	[gmt_modified] [datetime2](7) NULL
) ON [PRIMARY]
GO

