INSERT INTO `interaction` (`id`, `created_at`, `updated_at`, `biz_id`, `biz`, `read_cnt`, `like_cnt`, `collect_cnt`) VALUES (1, '1719213586574803', '1719213586574803', 181425015638462464, 'post', 50, 1, 1);

INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (1, '1719213586574803', '1719213586574803', 181425015638462464, '1111111111', '1111111111', '1111111111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (2, '1719213586574803', '1719213586574803', 181425032868663296, '111111111', '111111111', '111111111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (3, '1719213586574803', '1719213586574803', 181425046454013952, '11111111', '11111111', '11111111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (4, '1719213586574803', '1719213586574803', 181425061092134912, '1111111', '1111111', '1111111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (5, '1719213586574803', '1719213586574803', 181425073217867776, '111111', '111111', '111111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (6, '1719213586574803', '1719213586574803', 181425088577409024, '11111', '11111', '11111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (7, '1719213586574803', '1719213586574803', 181425102028541952, '1111', '1111', '1111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (8, '1719213586574803', '1719213586574803', 181425115697778688, '111', '111', '111', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (9, '1719213586574803', '1719213586574803', 181425138254745600, '11', '11', '11', 181424931144208384, 'publish');
INSERT INTO `post` (`id`, `created_at`, `updated_at`, `post_id`, `title`, `abstract`, `content`, `author_id`, `status`) VALUES (10, '1719213586574803', '1719213586574803', 181425152725094400, '1', '1', '1', 181424931144208384, 'publish');

INSERT INTO `user_collect` (`id`, `created_at`, `updated_at`, `user_id`, `biz_id`, `biz`, `collection_id`, `status`) VALUES (1, '1719213586574803', '1719213586574803', 181424931144208384, 181425015638462464, 'post', 0, 1);

INSERT INTO `user_like` (`id`, `created_at`, `updated_at`, `user_id`, `biz_id`, `biz`, `status`) VALUES (1, '1719213586574803', '1719213586574803', 181424931144208384, 181425015638462464, 'post', 1);

INSERT INTO `user` (`id`, `created_at`, `updated_at`, `user_id`, `user_name`, `email`, `password`) VALUES (1, '1719213586574803', '1719213586574803', 181424931144208384, 'root', '1227891082@qq.com', '$2a$10$09Wx5TZRcWyFJnSv8pr6xegAHG3mZE48eUkust9znC0JqNiK3YQeO');

INSERT INTO `comment` (`id`, `created_at`, `updated_at`, `comment_id`, `user_id`, `biz`, `biz_id`, `root_id`, `parent_id`, `content`) VALUES ('1', '1719213586574803', '1719213586574803', '182068922726486016', '181424931144208384', 'post', '181425015638462464', null, null, 'I am Root Comment');
INSERT INTO `comment` (`id`, `created_at`, `updated_at`, `comment_id`, `user_id`, `biz`, `biz_id`, `root_id`, `parent_id`, `content`) VALUES ('2', '1719213586574803', '1719213586574803', '182068981903921152', '181424931144208384', 'post', '181425015638462464', '1', '1', 'child comment 1');
INSERT INTO `comment` (`id`, `created_at`, `updated_at`, `comment_id`, `user_id`, `biz`, `biz_id`, `root_id`, `parent_id`, `content`) VALUES ('3', '1719213586574803', '1719213586574803', '182069018746687488', '181424931144208384', 'post', '181425015638462464', '1', '1', 'child comment 2');

INSERT INTO `file` (`id`, `created_at`, `updated_at`, `file_id`, `name`, `path`, `type`, `source`, `hash`) VALUES (1, '1719213586574803', '1719213586574803', 184200781623201792, '1716448088623493295.jpg', 'assets/file/1716448088623493295.jpg', 'image/jpeg', 'local', '747a20eb985c5dcbedde5612e506de03');
INSERT INTO `file` (`id`, `created_at`, `updated_at`, `file_id`, `name`, `path`, `type`, `source`, `hash`) VALUES (2, '1719213586574803', '1719213586574803', 184200781623201793, '1716448088624357809.jpg', 'assets/file/1716448088624357809.jpg', 'image/jpeg', 'local', '1e2267cf7526f877375ed7155bfd0f66');
INSERT INTO `file` (`id`, `created_at`, `updated_at`, `file_id`, `name`, `path`, `type`, `source`, `hash`) VALUES (3, '1719213586574803', '1719213586574803', 184200781623201794, '1716448088624300208.png', 'assets/file/1716448088624300208.png', 'image/png', 'local', '0f616e68388117f91b1f8e6ed678b807');

INSERT INTO `go-backend`.tag (id, created_at, updated_at, tag_id, user_id, tag_name) VALUES (1, '1719213586574803', '1719213586574803', 184614678431797248, 181424931144208384, 'Golang');
INSERT INTO `go-backend`.tag (id, created_at, updated_at, tag_id, user_id, tag_name) VALUES (2, '1719213586574803', '1719213586574803', 184614704025440256, 181424931144208384, 'Java');
INSERT INTO `go-backend`.tag (id, created_at, updated_at, tag_id, user_id, tag_name) VALUES (3, '1719213586574803', '1719213586574803', 184614731443605504, 181424931144208384, 'Php');
INSERT INTO `go-backend`.tag (id, created_at, updated_at, tag_id, user_id, tag_name) VALUES (4, '1719213586574803', '1719213586574803', 184614764402446336, 181424931144208384, 'Python');
INSERT INTO `go-backend`.tag (id, created_at, updated_at, tag_id, user_id, tag_name) VALUES (5, '1719213586574803', '1719213586574803', 184614779258671104, 181424931144208384, 'Rust');
