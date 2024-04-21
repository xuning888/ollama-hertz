create table prompt
(
    id          bigint(20)   not null auto_increment comment '主键',
    name        varchar(256) not null default '' comment '提示词的名称',
    description varchar(512) not null default '' comment '提示词的简单描述',
    temperature double       not null default 0.0 comment 'temperature',
    created_at  timestamp             default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at  timestamp             default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    deleted_at  timestamp    NULL     DEFAULT NULL COMMENT 'User information delete time',
    deleted     bigint(20)   not null default 0 comment '逻辑删除字段',
    primary key (id)
) default charset utf8mb4
  engine innodb comment '提示词主表';

create table prompt_item
(
    id         bigint(20)  not null auto_increment comment '主键',
    prompt_id  bigint(20)  not null default 0 comment '关联的主键',
    content    TEXT        not null comment 'content',
    role       varchar(32) not null default '' comment '角色',
    item_order int(16)     not null default 0 comment '提示词组织的顺序',
    created_at timestamp            default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp            default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    deleted_at timestamp   NULL     DEFAULT NULL COMMENT 'User information delete time',
    deleted    bigint(20)  not null default 0 comment '逻辑删除字段',
    primary key (id)
) default charset utf8mb4
  engine innodb comment '提示词细节';