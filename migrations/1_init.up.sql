CREATE TABLE IF NOT EXISTS subs(
    ID varchar(36) NOT NULL PRIMARY KEY,
    userID text NOT NULL,
    serviceName text NOT NULL,
    price text NOT NULL,
    startDate timestamp NOT NULL
);
