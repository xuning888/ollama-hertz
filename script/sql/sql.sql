create table yoyoyo_chat
(
    id         bigint auto_increment comment '主键',
    chat_id    varchar(64)  not null default '' comment '会话id',
    name       varchar(128) not null default '' comment '会话名称',
    file_refs  text         not null comment 'chatfile场景时用于存储附件的id, 格式是 [attach_id, attach_id], 会话关联附件',
    user_id    varchar(64)  not null default '' comment 'userId',
    chatbot_id varchar(64)  not null default '' comment '如果会话使用了chatbot就关联一下',
    created_at timestamp    not null default CURRENT_TIMESTAMP comment '创建时间',
    updated_at timestamp    not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    primary key (id),
    unique key unique_index_chat_id (chat_id),
    index idx_updated_at (updated_at)
) default charset = utf8mb4 comment '会话表';

create table yoyoyo_message
(
    id         bigint auto_increment comment '主键',
    msg_id     varchar(64)  not null default '消息id',
    content    text         not null comment '消息内容',
    role       varchar(64)  not null default '' comment '消息角色',
    group_id   varchar(128) not null default '' comment '分组号user和assistant的一次对话算一组',
    file_refs  text         not null comment 'chatfile场景时用于存储附件的id, 格式是 [attach_id, attach_id]，消息关联附件',
    chat_id    varchar(64)  not null default '' comment '会话id',
    in_memory  int          not null default 0 comment '使用用于上下文',
    created_at timestamp    not null default CURRENT_TIMESTAMP comment '创建时间',
    updated_at timestamp    not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    primary key (id),
    unique key unique_index_chat_id (msg_id)
) default charset = utf8mb4 comment '消息表';

create table yoyoyo_attachment
(
    id           bigint auto_increment comment '主键',
    attach_id    varchar(64)  not null default '' comment '附件id',
    name         varchar(128) not null default '' comment '附件名称',
    namespace    varchar(128) not null default '' comment '用于关联向量数据库的collection',
    download     varchar(500) not null default '' comment '下载附件的url',
    size         bigint       not null default 0 comment '附件大小',
    token_size   bigint       not null default 0 comment 'embedding所消耗的token数量',
    content_type varchar(64)  not null default '' comment '文件类型',
    md5          varchar(128) not null default '' comment 'md5,避免对附件重复的向量化',
    status       int          not null default 0 comment '附件处理状态',
    created_at   timestamp    not null default CURRENT_TIMESTAMP comment '创建时间',
    updated_at   timestamp    not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    primary key (id),
    unique key unique_index_attach_id (attach_id)
) default charset = utf8mb4
    comment '附件表, 存储了附件的名称和下载地址并关联了向量数据库的namespace, 附件不与会话强关联';

create table yoyoyo_chatbot
(
    id         bigint auto_increment comment '主键',
    chatbot_id varchar(64)  not null default '' comment '业务主键',
    name       varchar(128) not null default '' comment 'chatbotName',
    welcome    text         not null comment 'chatbot的介绍',
    created_at timestamp    not null default CURRENT_TIMESTAMP comment '创建时间',
    updated_at timestamp    not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    primary key (id),
    unique key unique_index_chatbot_id (chatbot_id)
) default charset = utf8mb4 comment 'chatbot';

create table yoyoyo_prompt
(
    id         bigint auto_increment comment '主键',
    prompt_id  varchar(64)  not null default '' comment '提示词id, 业务主键',
    role       varchar(64)  not null default '' comment '提示词的角色',
    prompt     text         not null comment '提示词或提示词模板',
    input      varchar(500) not null default '' comment '如果prompt是提示词模板, 那么这里应该存储inputValues, 用户输入时校验参数是否正确',
    chatbot_id varchar(64)  not null default '' comment '提示词关联的chatbot',
    sequence   int          not null default 0 comment '多个提示词组合为一个chatbot, chatbot使用提示词是应该有顺序, sequence就是顺序',
    created_at timestamp    not null default CURRENT_TIMESTAMP comment '创建时间',
    updated_at timestamp    not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    primary key (id),
    unique key unique_index_prompt_id (prompt_id)
) default charset = utf8mb4 comment '提示词';


create table yoyoyo_api_info
(
    id          bigint auto_increment comment '主键',
    endpoint    varchar(256) not null default '' comment '其实就是baseurl',
    token       varchar(256) not null default '' comment '认证用',
    model       varchar(128) not null default '' comment '模型名称, 模型id',
    description varchar(256) not null default '' comment '备注和说明',
    type        int          not null default 0 comment 'API类型: 0-OpenAI, 1-Azure OpenAI 以此类推',
    llm_type    int          not null default 0 comment '大模型类型: 0-generate, 1-embedding, 2-both',
    deleted     int          not null default 0 comment '逻辑删除',
    created_at  timestamp    not null default CURRENT_TIMESTAMP comment '创建时间',
    updated_at  timestamp    not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
    primary key (id)
) default charset = utf8mb4 comment 'api管理';



