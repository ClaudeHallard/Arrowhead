CREATE TABLE Services (
    id                 INTEGER      PRIMARY KEY,
    serviceDefinition  CHAR (50),
    SystemName         TEXT         NOT NULL,
    metaData           CHAR (45)    DEFAULT 0,
    Port               INTEGER (10),
    authenticationInfo CHAR (50),
    serviceURI         CHAR (100),
    endOfValidity      CHAR (50),
    secure             CHAR (50),
    address            CHAR (100),
    version            CHAR (50),
    interfaces         CHAR (100) 
);