--
-- File generated with SQLiteStudio v3.3.3 on fre jan 28 13:48:11 2022
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Interfaces
CREATE TABLE Interfaces (
    idInterface   INTEGER,
    serviceID     INTEGER,
    interfaceName CHAR,
    createdAt     CHAR,
    updatedAt     CHAR,
    PRIMARY KEY (
        idInterface
    )
    AUTOINCREMENT,
    FOREIGN KEY (
        serviceID
    )
    REFERENCES Services (id) ON DELETE CASCADE
);


-- Table: MetaData
CREATE TABLE MetaData (
    idMetaData INTEGER,
    serviceID  INTEGER,
    metaData   CHAR,
    PRIMARY KEY (
        idMetaData
    )
    AUTOINCREMENT,
    FOREIGN KEY (
        serviceID
    )
    REFERENCES Services (id) ON DELETE CASCADE
);


-- Table: Services
CREATE TABLE Services (
    id                         INTEGER NOT NULL,
    serviceDefinition          CHAR    NOT NULL,
    serviceDefinitionCreatedAt CHAR    NOT NULL
                                       DEFAULT 'defaultValue',
    serviceDefinitionUpdatedAt CHAR    NOT NULL
                                       DEFAULT 'defaultValue',
    systemName                 CHAR    NOT NULL,
    address                    CHAR    NOT NULL,
    port                       INTEGER NOT NULL,
    authenticationInfo         CHAR    NOT NULL,
    providerCreatedAt          CHAR    NOT NULL
                                       DEFAULT 'defaultValue',
    providerUpdatedAt          CHAR    NOT NULL
                                       DEFAULT 'defaultValue',
    serviceURI                 CHAR    NOT NULL,
    endOfValidity              CHAR    NOT NULL,
    secure                     CHAR    NOT NULL,
    version                    CHAR    NOT NULL,
    createdAt                  CHAR    NOT NULL
                                       DEFAULT 'defaultValue',
    updatedAt                  CHAR    NOT NULL
                                       DEFAULT 'defaultValue',
    PRIMARY KEY (
        id
    )
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;