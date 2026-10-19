CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,

    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    city VARCHAR(100) NOT NULL,
    birthday DATE NOT NULL,
    profile_picture VARCHAR(500) NULL,
    description TEXT NULL,

    language ENUM('en', 'al') NOT NULL DEFAULT 'en',
    theme ENUM('dark', 'light') NOT NULL DEFAULT 'light',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    strikes_count INT UNSIGNED NOT NULL DEFAULT 0,

    INDEX idx_users_city (city)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;

create table interests (
    id bigint unsigned auto_increment primary key,
    name varchar(100) not null unique
) engine=InnoDB;

create table user_interests (
    user_id bigint unsigned not null,
    interest_id bigint unsigned not null,
    primary key (user_id, interest_id),
    foreign key (user_id)
    references users(id) on delete cascade,
    foreign key (interest_id)
    references interests(id) on delete cascade
) engine=InnoDB;

insert into interests (name) values ('Music'), ('Hiking'), ('Gaming'), ('Travel'), ('Photography'), ('Board Games'), ('Sports'), ('Art'), ('Food'), ('Nightlife');

create table user_connections (
    user_id bigint unsigned not null,
    connected_user_id bigint unsigned not null,
    created_at timestamp not null default current_timestamp,
    primary key (user_id, connected_user_id),
    foreign key (user_id) 
    references users(id) on delete cascade,
    foreign key (connected_user_id)
    references users(id) on delete cascade,
    check (user_id <> connected_user_id)
) engine=InnoDB;

create table blocked_users (
    user_id bigint unsigned not null,
    blocked_user_id bigint unsigned not null,
    created_at timestamp not null default current_timestamp,
    primary key (user_id, blocked_user_id),
    foreign key (user_id)
    references users(id) on delete cascade,
    foreign key (blocked_user_id)
    references users(id) on delete cascade,
    check (user_id <> blocked_user_id)
)engine=InnoDB;

create table emergency_contacts(
    id bigint unsigned auto_increment primary key,
    user_id bigint unsigned not null,
    name varchar(200) not null,
    phone varchar(30) not null,
    created_at timestamp not null default current_timestamp,
    foreign key (user_id)
    references users(id) on delete cascade,
    index idx_emergency_contacts_user (user_id)
) engine=InnoDB;

create table categories (
    id bigint unsigned auto_increment primary key,
    name varchar(100) not null unique,
    description varchar(255) null,
    image_url varchar(500) null
) engine=InnoDB;

CREATE TABLE events (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    creator_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,

    title VARCHAR(200) NOT NULL,
    description TEXT NULL,
    image_url VARCHAR(500) NULL,

    location_name VARCHAR(200) NOT NULL,
city VARCHAR(100) NOT NULL,
address VARCHAR(500) NULL,
    latitude DECIMAL(10, 7) NULL,
    longitude DECIMAL(10, 7) NULL,

    starts_at DATETIME NOT NULL,
    ends_at DATETIME NOT NULL,
    registration_deadline DATETIME NULL,

    min_participants INT UNSIGNED NULL,
    max_participants INT UNSIGNED NULL,

    price_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    currency CHAR(3) NOT NULL DEFAULT 'ALL',

    status ENUM(
        'draft',
        'published',
        'cancelled'
    ) NOT NULL DEFAULT 'published',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (creator_id)
        REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (category_id)
        REFERENCES categories(id) ON DELETE RESTRICT,

    CHECK (ends_at > starts_at),
    CHECK (
        min_participants IS NULL
        OR min_participants >= 1
    ),
    CHECK (
        max_participants IS NULL
        OR max_participants >= 1
    ),
    CHECK (
        min_participants IS NULL
        OR max_participants IS NULL
        OR min_participants <= max_participants
    ),
    CHECK (price_amount >= 0),

    INDEX idx_events_discovery (status, starts_at, city),
    INDEX idx_events_creator (creator_id),
    INDEX idx_events_category (category_id)
) ENGINE=InnoDB;

CREATE TABLE event_interests (
    event_id BIGINT UNSIGNED NOT NULL,
    interest_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (event_id, interest_id),

    FOREIGN KEY (event_id)
        REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (interest_id)
        REFERENCES interests(id) ON DELETE RESTRICT
) ENGINE=InnoDB;

CREATE TABLE app_groups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    creator_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,

    title VARCHAR(200) NOT NULL,
    description TEXT NULL,
    image_url VARCHAR(500) NULL,
    city VARCHAR(100) NOT NULL,

    status ENUM('active', 'archived')
        NOT NULL DEFAULT 'active',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (creator_id)
        REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (category_id)
        REFERENCES categories(id) ON DELETE RESTRICT,

    INDEX idx_groups_city (city),
    INDEX idx_groups_creator (creator_id)
) ENGINE=InnoDB;

CREATE TABLE group_interests (
    group_id BIGINT UNSIGNED NOT NULL,
    interest_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (group_id, interest_id),

    FOREIGN KEY (group_id)
        REFERENCES app_groups(id) on delete cascade,
FOREIGN KEY (interest_id)
        REFERENCES interests(id) ON DELETE RESTRICT
) ENGINE=InnoDB;

create table group_members (
    group_id bigint unsigned not null,
    user_id bigint unsigned not null,
    role enum ('owner', 'admin', 'member')
    not null default 'member',
    joined_at timestamp not null default current_timestamp,
    primary key (group_id, user_id),
    foreign key(group_id)
    references app_groups(id) on delete cascade,
    foreign key (user_id)
    references users(id) on delete cascade,
    index idx_group_members_user (user_id)
)engine=InnoDB;

create table group_events (
    group_id bigint unsigned not null,
    event_id bigint unsigned not null,
    linked_at timestamp not null default current_timestamp,
    primary key (group_id, event_id),
    foreign key (group_id)
    references app_groups(id) on delete cascade,
    foreign key (event_id) 
    references events(id) on delete cascade
) engine=InnoDB;

CREATE TABLE chats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    type ENUM('direct', 'event', 'group') NOT NULL,

    event_id BIGINT UNSIGNED NULL,
    group_id BIGINT UNSIGNED NULL,

    user_one_id BIGINT UNSIGNED NULL,
    user_two_id BIGINT UNSIGNED NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_chat_event (event_id),
    UNIQUE KEY uq_chat_group (group_id),
    UNIQUE KEY uq_chat_direct_pair (user_one_id, user_two_id),

    FOREIGN KEY (event_id)
        REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id)
        REFERENCES app_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_one_id)
        REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_two_id)
        REFERENCES users(id) ON DELETE CASCADE,
        

    CHECK (
        (type = 'direct'
            AND event_id IS NULL
            AND group_id IS NULL
            AND user_one_id IS NOT NULL
            AND user_two_id IS NOT NULL
            AND user_one_id < user_two_id)
        OR
        (type = 'event'
            AND event_id IS NOT NULL
            AND group_id IS NULL
            AND user_one_id IS NULL
            AND user_two_id IS NULL)
        OR
        (type = 'group'
            AND group_id IS NOT NULL
            AND event_id IS NULL
            AND user_one_id IS NULL
            AND user_two_id IS NULL)
    )
) ENGINE=InnoDB;

CREATE TABLE chat_members (
    chat_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,

    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at DATETIME NULL,
    last_read_at DATETIME NULL,

    PRIMARY KEY (chat_id, user_id),

    FOREIGN KEY (chat_id)
        REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_chat_members_user (user_id, left_at)
) ENGINE=InnoDB;

CREATE TABLE messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    chat_id BIGINT UNSIGNED NOT NULL,
    sender_id BIGINT UNSIGNED NOT NULL,

    message_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME NULL,
    deleted_at DATETIME NULL,

    FOREIGN KEY (chat_id)
        REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id)
        REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_messages_chat_time (chat_id, created_at)
) ENGINE=InnoDB;

CREATE TABLE notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    recipient_id BIGINT UNSIGNED NOT NULL,
    actor_id BIGINT UNSIGNED NULL,

    type ENUM(
        'connection_request',
        'connection_accepted',
        'event_invitation',
        'event_update',
        'event_cancelled',
        'event_reminder',
        'group_invitation',
        'group_update',
        'new_message',
        'event_chat_message',
        'group_chat_message',
        'user_report_update',
        'strike_received'
    ) NOT NULL,

    entity_type ENUM(
        'user',
        'event',
        'group',
        'chat',
        'report',
        'strike'
    ) NULL,

    entity_id BIGINT UNSIGNED NULL,

    title VARCHAR(200) NOT NULL,
    body TEXT NULL,

    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at DATETIME NULL,

    FOREIGN KEY (recipient_id)
        REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id)
        REFERENCES users(id) ON DELETE SET NULL,

    INDEX idx_notifications_recipient (recipient_id, is_read, created_at)
) ENGINE=InnoDB;

CREATE TABLE connections (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    requester_id BIGINT UNSIGNED NOT NULL,
    recipient_id BIGINT UNSIGNED NOT NULL,

    status ENUM('pending', 'accepted', 'rejected')
        NOT NULL DEFAULT 'pending',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uq_connection_pair (requester_id, recipient_id),

    FOREIGN KEY (requester_id)
        REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id)
        REFERENCES users(id) ON DELETE CASCADE,

    CHECK (requester_id <> recipient_id),

    INDEX idx_connections_recipient (recipient_id, status)
) ENGINE=InnoDB;

CREATE TABLE event_participants (
    event_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,

    status ENUM(
        'interested',
        'going',
        'attended',
        'cancelled'
    ) NOT NULL DEFAULT 'interested',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (event_id, user_id),

    FOREIGN KEY (event_id)
        REFERENCES events(id) ON DELETE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_event_participants_user (user_id, status)
) ENGINE=InnoDB;

CREATE TABLE user_reports (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    reporter_id BIGINT UNSIGNED NOT NULL,
    reported_user_id BIGINT UNSIGNED NOT NULL,

    reason VARCHAR(100) NOT NULL,
    details TEXT NULL,

    status ENUM(
        'pending',
        'reviewed',
        'dismissed',
        'action_taken'
    ) NOT NULL DEFAULT 'pending',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at DATETIME NULL,

    FOREIGN KEY (reporter_id)
        REFERENCES users(id) ON DELETE RESTRICT,

    FOREIGN KEY (reported_user_id)
        REFERENCES users(id) ON DELETE RESTRICT,

    CHECK (reporter_id <> reported_user_id),

    INDEX idx_reports_reported_status (reported_user_id, status)
) ENGINE=InnoDB;

CREATE TABLE user_strikes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    event_id BIGINT UNSIGNED NULL,
    issued_by BIGINT UNSIGNED NOT NULL,

    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE RESTRICT,

    FOREIGN KEY (event_id)
        REFERENCES events(id) ON DELETE SET NULL,

    FOREIGN KEY (issued_by)
        REFERENCES users(id) ON DELETE RESTRICT,

    INDEX idx_strikes_user (user_id)
) ENGINE=InnoDB;