USE [kequ5060]
GO

/****** Object:  Table [dbo].[guild_pblock]    Script Date: 2023/8/12 19:39:19 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[guild_pblock](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[guild_id] [varchar](32) NOT NULL,
	[user_id] [varchar](32) NOT NULL,
	[admin_id] [varchar](32) NULL,
	[gmt_modified] [datetime2](7) NULL,
	[ispblock] [bit] NULL,
 CONSTRAINT [PK_guild_pblock] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

