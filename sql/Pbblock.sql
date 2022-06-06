USE [kequ5060]
GO

/****** Object:  Table [dbo].[zbot_pblock]    Script Date: 2022/6/6 10:58:40 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[zbot_pblock](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[user_id] [bigint] NOT NULL,
	[admin_id] [bigint] NULL,
	[gmt_modified] [datetime2](7) NULL,
	[ispblock] [bit] NULL,
 CONSTRAINT [PK_zbot_pblock] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

