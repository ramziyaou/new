CREATE TABLE `users`
(
    id   int,
    username text,
    password text,
    wallets text
);

CREATE TABLE `crypto_wallets`
(
    username text,
    name text,
    amount bigint
);


CREATE TABLE `start_stop_checks`
(
    username text,
    name text,
    stop int,
    start int
);

INSERT INTO `users` (`id`, `username`, `password`, `wallets`)
VALUES (1, 'Ramziya', '1234', ''),
       (2, 'two', 'twopw', 'tone '),
       (3, 'three', 'threepw', 'thone thtwo '),
       (4, 'four', 'fourpw', 'fone ftwo fthree '),
       (6, 'six', 'sixpw', 'sone ');


INSERT INTO `crypto_wallets` (`username`, `name`, `amount`)
VALUES ('two', 'tone', 0),
       ('three', 'thone', 0),
       ('three', 'thtwo', 0),
       ('four', 'fone', 0),
       ('four', 'ftwo', 0),
       ('four', 'fthree', 0),
       ('six', 'sone', 0);


INSERT INTO `start_stop_checks` (`username`, `name`, `stop`, `start`)
VALUES ('two', 'tone', 0,0),
       ('three', 'thone', 0, 0),
       ('three', 'thtwo', 0, 0),
       ('four', 'fone', 0,0),
       ('four', 'ftwo', 0,0),
       ('four', 'fthree', 0,0),
       ('six', 'sone', 0,0);