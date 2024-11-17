DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS tags;

CREATE TABLE authors (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         name VARCHAR(100) NOT NULL,
                         email VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE posts (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       title VARCHAR(200) NOT NULL,
                       content TEXT NOT NULL,
                       author_id INT,
                       CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES authors (id)
                           ON UPDATE CASCADE
                           ON DELETE SET NULL
);

CREATE TABLE tags (
                      id INT AUTO_INCREMENT PRIMARY KEY,
                      name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE post_tags (
                           post_id INT,
                           tag_id INT,
                           PRIMARY KEY (post_id, tag_id),
                           CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts (id)
                               ON DELETE CASCADE,
                           CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id)
                               ON DELETE CASCADE
);



INSERT INTO authors (name, email)
VALUES ('Alice Dupont', 'alice.dupont@example.com'),
       ('Bob Martin', 'bob.martin@example.com'),
       ('Clara Petit', 'clara.petit@example.com'),
       ('David Moreau', 'david.moreau@example.com'),
       ('Emma Lefevre', 'emma.lefevre@example.com'),
       ('François Bernard', 'francois.bernard@example.com'),
       ('Gina Rousseau', 'gina.rousseau@example.com'),
       ('Hugo Caron', 'hugo.caron@example.com'),
       ('Isabelle Leroy', 'isabelle.leroy@example.com'),
       ('Julien Simon', 'julien.simon@example.com');

INSERT INTO posts (title, content, author_id)
VALUES ('Post 1', 'Contenu du post 1', 1),
       ('Post 2', 'Contenu du post 2', 2),
       ('Post 3', 'Contenu du post 3', 3),
       ('Post 4', 'Contenu du post 4', 4),
       ('Post 5', 'Contenu du post 5', 5),
       ('Post 6', 'Contenu du post 6', 1),
       ('Post 7', 'Contenu du post 7', 2),
       ('Post 8', 'Contenu du post 8', 3),
       ('Post 9', 'Contenu du post 9', 4),
       ('Post 10', 'Contenu du post 10', 5),
       ('Post 11', 'Contenu du post 11', 1),
       ('Post 12', 'Contenu du post 12', 2),
       ('Post 13', 'Contenu du post 13', 3),
       ('Post 14', 'Contenu du post 14', 4),
       ('Post 15', 'Contenu du post 15', 5),
       ('Post 16', 'Contenu du post 16', 6),
       ('Post 17', 'Contenu du post 17', 7),
       ('Post 18', 'Contenu du post 18', 8),
       ('Post 19', 'Contenu du post 19', 9),
       ('Post 20', 'Contenu du post 20', 10),
       ('Post 21', 'Contenu du post 21', 1),
       ('Post 22', 'Contenu du post 22', 2),
       ('Post 23', 'Contenu du post 23', 3),
       ('Post 24', 'Contenu du post 24', 4),
       ('Post 25', 'Contenu du post 25', 5),
       ('Post 26', 'Contenu du post 26', 6),
       ('Post 27', 'Contenu du post 27', 7),
       ('Post 28', 'Contenu du post 28', 8),
       ('Post 29', 'Contenu du post 29', 9),
       ('Post 30', 'Contenu du post 30', 10),
       ('Post 31', 'Contenu du post 31', 1),
       ('Post 32', 'Contenu du post 32', 2),
       ('Post 33', 'Contenu du post 33', 3),
       ('Post 34', 'Contenu du post 34', 4),
       ('Post 35', 'Contenu du post 35', 5),
       ('Post 36', 'Contenu du post 36', 6),
       ('Post 37', 'Contenu du post 37', 7),
       ('Post 38', 'Contenu du post 38', 8),
       ('Post 39', 'Contenu du post 39', 9),
       ('Post 40', 'Contenu du post 40', 10),
       ('Post 41', 'Contenu du post 41', 1),
       ('Post 42', 'Contenu du post 42', 2),
       ('Post 43', 'Contenu du post 43', 3),
       ('Post 44', 'Contenu du post 44', 4),
       ('Post 45', 'Contenu du post 45', 5),
       ('Post 46', 'Contenu du post 46', 6),
       ('Post 47', 'Contenu du post 47', 7),
       ('Post 48', 'Contenu du post 48', 8),
       ('Post 49', 'Contenu du post 49', 9),
       ('Post 50', 'Contenu du post 50', 10);

INSERT INTO tags (name)
VALUES ('Tech'),
       ('Science'),
       ('Art'),
       ('Travel'),
       ('Food'),
       ('Health'),
       ('Education'),
       ('Lifestyle');

INSERT INTO post_tags (post_id, tag_id)
VALUES (1, 1),
       (1, 2),
       (2, 3),
       (3, 1),
       (3, 4),
       (4, 5),
       (5, 6),
       (6, 1),
       (7, 7),
       (8, 8),
       (9, 5),
       (10, 6);
