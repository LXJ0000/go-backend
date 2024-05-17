INSERT INTO `interaction` (`id`, `created_at`, `updated_at`, `biz_id`, `biz`, `read_cnt`, `like_cnt`, `collect_cnt`) VALUES (1, '2024-05-15 23:22:42.485', '2024-05-15 23:22:42.485', 181425015638462464, 'post', 50, 1, 1);

INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (1, '2024-05-15 23:17:56.143', '2024-05-15 23:17:56.143', 181425015638462464, '1111111111', '1111111111', '1111111111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (2, '2024-05-15 23:18:00.252', '2024-05-15 23:18:00.252', 181425032868663296, '111111111', '111111111', '111111111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (3, '2024-05-15 23:18:03.491', '2024-05-15 23:18:03.491', 181425046454013952, '11111111', '11111111', '11111111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (4, '2024-05-15 23:18:06.980', '2024-05-15 23:18:06.980', 181425061092134912, '1111111', '1111111', '1111111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (5, '2024-05-15 23:18:09.872', '2024-05-15 23:18:09.872', 181425073217867776, '111111', '111111', '111111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (6, '2024-05-15 23:18:13.533', '2024-05-15 23:18:13.533', 181425088577409024, '11111', '11111', '11111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (7, '2024-05-15 23:18:16.740', '2024-05-15 23:18:16.740', 181425102028541952, '1111', '1111', '1111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (8, '2024-05-15 23:18:18.462', '2024-05-15 23:18:18.462', 181425115697778688, '111', '111', '111', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (9, '2024-05-15 23:18:23.841', '2024-05-15 23:18:23.841', 181425138254745600, '11', '11', '11', 181424931144208384, 2);
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (10, '2024-05-15 23:18:27.290', '2024-05-15 23:18:27.290', 181425152725094400, '1', '1', '1', 181424931144208384, 2);

INSERT INTO `user_collect` (`id`, `created_at`, `updated_at`, `user_id`, `biz_id`, `biz`, `collection_id`, `status`) VALUES (1, '2024-05-15 23:22:46.224', '2024-05-15 23:22:46.224', 181424931144208384, 181425015638462464, 'post', 0, 1);

INSERT INTO `user_like` (`id`, `created_at`, `updated_at`, `user_id`, `biz_id`, `biz`, `status`) VALUES (1, '2024-05-15 23:22:42.486', '2024-05-15 23:22:42.486', 181424931144208384, 181425015638462464, 'post', 1);

INSERT INTO `user` (`id`, `created_at`, `updated_at`, `user_id`, `user_name`, `email`, `password`) VALUES (1, '2024-05-15 23:17:37.546', '2024-05-15 23:17:37.546', 181424931144208384, 'root', '1227891082@qq.com', '$2a$10$09Wx5TZRcWyFJnSv8pr6xegAHG3mZE48eUkust9znC0JqNiK3YQeO');

INSERT INTO `comment` (`id`, `created_at`, `updated_at`, `comment_id`, `user_id`, `biz`, `biz_id`, `root_id`, `parent_id`, `content`) VALUES ('1', '2024-05-17 17:56:53.854', '2024-05-17 17:56:53.854', '182068922726486016', '181424931144208384', 'post', '181425015638462464', null, null, 'I am Root Comment');
INSERT INTO `comment` (`id`, `created_at`, `updated_at`, `comment_id`, `user_id`, `biz`, `biz_id`, `root_id`, `parent_id`, `content`) VALUES ('2', '2024-05-17 17:57:07.964', '2024-05-17 17:57:07.964', '182068981903921152', '181424931144208384', 'post', '181425015638462464', '1', '1', 'child comment 1');
INSERT INTO `comment` (`id`, `created_at`, `updated_at`, `comment_id`, `user_id`, `biz`, `biz_id`, `root_id`, `parent_id`, `content`) VALUES ('3', '2024-05-17 17:57:16.747', '2024-05-17 17:57:16.747', '182069018746687488', '181424931144208384', 'post', '181425015638462464', '1', '1', 'child comment 2');
