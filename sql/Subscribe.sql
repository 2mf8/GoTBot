USE [kequ5060]
GO

/****** Object:  Table [dbo].[zbot_replace]    Script Date: 2022/6/6 11:00:58 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[zbot_replace](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[orgin_group_id] [bigint] NOT NULL,
	[replace_group_id] [bigint] NOT NULL,
	[admin_id] [bigint] NOT NULL,
	[is_true] [bit] NULL,
 CONSTRAINT [PK_zbot_replace] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [dbo].[zbot_replace] ADD  CONSTRAINT [DF_zbot_replace_o_group_id]  DEFAULT ((0)) FOR [orgin_group_id]
GO

ALTER TABLE [dbo].[zbot_replace] ADD  CONSTRAINT [DF_zbot_replace_r_group_id]  DEFAULT ((0)) FOR [replace_group_id]
GO

ALTER TABLE [dbo].[zbot_replace] ADD  CONSTRAINT [DF_zbot_replace_admin_id]  DEFAULT ((0)) FOR [admin_id]
GO

