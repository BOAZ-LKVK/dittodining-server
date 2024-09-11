CREATE DATABASE IF NOT EXISTS `dittodining` DEFAULT CHARACTER SET utf8mb4;

USE `dittodining`;

CREATE TABLE `restaurant`
(
    `restaurant_id`            BIGINT         NOT NULL AUTO_INCREMENT COMMENT '음식점 ID',
    `name`                     VARCHAR(255)   NOT NULL COMMENT '음식점 이름',
    `address`                  VARCHAR(1024)  NOT NULL COMMENT '음식점 주소',
    `description`              TEXT           NOT NULL COMMENT '음식점 한줄 소개',
    `maximum_price_per_person` DECIMAL(11, 2) NOT NULL COMMENT '1인당 최대 가격대',
    `minimum_price_per_person` DECIMAL(11, 2) NOT NULL COMMENT '1인당 최소 가격대',
    `latitude`                 DECIMAL(11, 8) NOT NULL COMMENT '위도',            -- TODO: 위치 기반 스캔을 위해 SPATIAL INDEX가 필요할 수도 있음
    `longitude`                DECIMAL(11, 8) NOT NULL COMMENT '경도',            -- TODO: 위치 기반 스캔을 위해 SPATIAL INDEX가 필요할 수도 있음
    /* business_hours_json 예시
    [
	    {
            "dayOfWeek": "DAY_OF_WEEK_SUNDAY",
            "openTime": "10:00",
            "closingTime": "20:00",
            "isClosedDay": false
        },
        {
            "dayOfWeek": "DAY_OF_WEEK_MONDAY",
            "openTime": "10:00",
            "closingTime": "20:00",
            "isClosedDay": false
        },
    ]
     */
    `business_hours_json`      JSON           NOT NULL COMMENT '영업 시간 목록 JSON', -- BusinessHour 객체의 JSON
    `recommendation_score`     DECIMAL(5, 2)  NOT NULL COMMENT '추천 점수',
    `average_score_from_naver` DECIMAL(5, 2)  NULL COMMENT '네이버 평균 평점',
    `average_score_from_kakao` DECIMAL(5, 2)  NULL COMMENT '카카오 평균 평점',
    `created_at`               DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`               DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`restaurant_id`),
    KEY `idx_restaurant_m1` (`created_at`),
    KEY `idx_restaurant_m2` (`updated_at`),
    KEY `idx_restaurant_m3` (`latitude`, `longitude`, `recommendation_score`),
    KEY `idx_restaurant_m4` (`recommendation_score`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '음식점';

CREATE TABLE `restaurant_image`
(
    `restaurant_image_id` BIGINT        NOT NULL AUTO_INCREMENT COMMENT '음식점 이미지 ID',
    `restaurant_id`       BIGINT        NOT NULL COMMENT '음식점 ID',
    `image_url`           varchar(1024) NOT NULL COMMENT '이미지 URL',
    `created_at`          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`restaurant_image_id`),
    KEY `idx_restaurant_image_m1` (`created_at`),
    KEY `idx_restaurant_image_m2` (`updated_at`),
    KEY `idx_restaurant_image_m3` (`restaurant_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '음식점 이미지';


CREATE TABLE `restaurant_menu`
(
    `restaurant_menu_id` BIGINT         NOT NULL AUTO_INCREMENT COMMENT '음식점 메뉴 ID',
    `restaurant_id`      BIGINT         NOT NULL COMMENT '음식점 ID',
    `name`               varchar(255)   NOT NULL COMMENt '메뉴 이름',
    `price`              decimal(11, 2) NOT NULL COMMENT '가격',
    `description`        TEXT           NULL COMMENT '메뉴 설명',
    `image_url`          varchar(1024)  NULL COMMENT '이미지 URL',
    `created_at`         DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`         DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`restaurant_menu_id`),
    KEY `idx_restaurant_menu_m1` (`created_at`),
    KEY `idx_restaurant_menu_m2` (`updated_at`),
    KEY `idx_restaurant_menu_m3` (`restaurant_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '음식점 메뉴';

CREATE TABLE `restaurant_review`
(
    `restaurant_review_id` BIGINT        NOT NULL AUTO_INCREMENT COMMENT '음식점 리뷰 ID',
    `restaurant_id`        BIGINT        NOT NULL COMMENT '음식점 ID',
    `writer_name`          VARCHAR(255)  NOT NULL COMMENT '작성자 이름',
    `score`                DECIMAL(5, 2) NULL COMMENT '평점',
    `content`              TEXT          NULL COMMENT '내용',
    `wrote_at`             DATETIME      NOT NULL COMMENT '작성 일시',
    `created_at`           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`restaurant_review_id`),
    KEY `idx_restaurant_review_m1` (`created_at`),
    KEY `idx_restaurant_review_m2` (`updated_at`),
    KEY `idx_restaurant_review_m3` (`restaurant_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '음식점 리뷰';


CREATE TABLE `restaurant_recommendation_request`
(
    `restaurant_recommendation_request_id` BIGINT         NOT NULL AUTO_INCREMENT COMMENT '음식점 추천 요청 ID',
    `user_id`                              BIGINT         NULL COMMENT '사용자 ID',
    `latitude`                             DECIMAL(11, 8) NOT NULL COMMENT '유저 위도',
    `longitude`                            DECIMAL(11, 8) NOT NULL COMMENT '유저 경도',
    `requested_at`                         DATETIME       NOT NULL COMMENT '요청 일시',
    `created_at`                           DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`                           DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`restaurant_recommendation_request_id`),
    KEY `idx_restaurant_recommendation_request_m1` (`created_at`),
    KEY `idx_restaurant_recommendation_request_m2` (`updated_at`),
    KEY `idx_restaurant_recommendation_request_m3` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '음식점 추천 요청';


CREATE TABLE `restaurant_recommendation`
(
    `restaurant_recommendation_id`         BIGINT         NOT NULL AUTO_INCREMENT COMMENT '음식점 추천 ID',
    `restaurant_recommendation_request_id` BIGINT         NOT NULL COMMENT '음식점 추천 요청 ID',
    `restaurant_id`                        BIGINT         NOT NULL COMMENT '음식점 ID',
    `distance_in_meters`                   DECIMAL(11, 2) NOT NULL COMMENT '미터 단위 거리',
    `created_at`                           DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`                           DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`restaurant_recommendation_id`),
    KEY `idx_restaurant_recommendation_m1` (`created_at`),
    KEY `idx_restaurant_recommendation_m2` (`updated_at`),
    KEY `idx_restaurant_recommendation_m3` (`restaurant_recommendation_request_id`),
    KEY `idx_restaurant_recommendation_m4` (`restaurant_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '음식점 추천';


CREATE TABLE `selected_restaurant_recommendation`
(
    `selected_restaurant_recommendation_id` BIGINT   NOT NULL AUTO_INCREMENT COMMENT '선택된 음식점 추천 ID',
    `restaurant_recommendation_request_id`  BIGINT   NOT NULL COMMENT '음식점 추천 요청 ID',
    `restaurant_recommendation_id`          BIGINT   NOT NULL COMMENT '음식점 추천 ID',
    `restaurant_id`                         BIGINT   NOT NULL COMMENT '음식점 ID',
    `created_at`                            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성 일시',
    `updated_at`                            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정 일시',
    PRIMARY KEY (`selected_restaurant_recommendation_id`),
    KEY `idx_user_selected_restaurant_recommendation_m1` (`created_at`),
    KEY `idx_user_selected_restaurant_recommendation_m2` (`updated_at`),
    KEY `idx_user_selected_restaurant_recommendation_m3` (`restaurant_recommendation_request_id`),
    KEY `idx_user_selected_restaurant_recommendation_m4` (`restaurant_recommendation_id`),
    KEY `idx_user_selected_restaurant_recommendation_m5` (`restaurant_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4` COMMENT '선택된 음식점 추천';
