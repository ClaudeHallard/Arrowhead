--
-- File generated with SQLiteStudio v3.3.3 on ons jan 19 15:02:40 2022
--
-- Text encoding used: System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Interfaces
CREATE TABLE Interfaces (
    idInterfaces  INTEGER PRIMARY KEY,
    serviceID     INTEGER REFERENCES Services (id) ON DELETE CASCADE,
    interfaceName CHAR,
    createdAt     CHAR,
    updatedAt     CHAR
);


-- Table: MetaData
CREATE TABLE MetaData (
    idMetaData INTEGER PRIMARY KEY,
    ServiceID  INTEGER REFERENCES Services (id) ON DELETE CASCADE,
    MetaData   CHAR
);


-- Table: Services
CREATE TABLE Services (
    id                         INTEGER PRIMARY KEY,
    serviceDefinition          CHAR,
    ServiceDefinitionCreatedAt CHAR,
    ServiceDefinitionUpdatedAt CHAR,
    SystemName                 CHAR,
    Address                    CHAR,
    Port                       INTEGER,
    authenticationInfo         CHAR,
    providerCreatedAt          CHAR,
    providerUpdatedAt          CHAR,
    serviceURI                 CHAR,
    endOfValidity              CHAR,
    secure                     CHAR,
    version                    CHAR,
    createdAt                  CHAR,
    updatedAt                  CHAR
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
