USE [kequ5060]
GO

/****** Object:  Table [dbo].[zbot_plugin_switch]    Script Date: 2021/10/30 6:20:33 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[zbot_plugin_switch](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[group_id] [bigint] NOT NULL,
	[plugin_name] [varchar](50) NOT NULL,
	[gmt_modified] [datetime2](7) NULL,
	[stop] [bit] NULL,
 CONSTRAINT [PK_zbot_plugin_switch] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO


