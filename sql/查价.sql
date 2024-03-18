USE [kequ5060]
GO

/****** Object:  Table [dbo].[guild_price]    Script Date: 2023/8/12 19:38:57 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[guild_price](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[guild_id] [varchar](32) NOT NULL,
	[channel_id] [varchar](32) NOT NULL,
	[brand] [nvarchar](50) NULL,
	[item] [nvarchar](100) NOT NULL,
	[price] [nvarchar](100) NULL,
	[shipping] [nvarchar](100) NULL,
	[updater] [nvarchar](32) NULL,
	[gmt_modified] [datetime2](7) NULL,
 CONSTRAINT [PK_guild_price] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

