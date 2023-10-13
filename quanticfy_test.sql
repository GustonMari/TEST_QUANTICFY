-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.

CREATE TABLE `ChannelType` (
    `ChannelTypeID` smallint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `Name` varchar(30)  NOT NULL ,
    PRIMARY KEY (
        `ChannelTypeID`
    )
);

CREATE TABLE `EventType` (
    `EventTypeID` smallint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `Name` varchar(30)  NOT NULL ,
    PRIMARY KEY (
        `EventTypeID`
    )
);

CREATE TABLE `Content` (
    `ContentID` int UNSIGNED AUTO_INCREMENT NOT NULL ,
    `ClientContentID` bigint UNSIGNED NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `ContentID`
    )
);

CREATE TABLE `Customer` (
    `CustomerID` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `ClientCustomerID` bigint UNSIGNED NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `CustomerID`
    )
);

CREATE TABLE `CustomerData` (
    `CustomerChannelID` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `CustomerID` bigint UNSIGNED NOT NULL ,
    `ChannelTypeID` smallint UNSIGNED NOT NULL ,
    `ChannelValue` varchar(600)  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `CustomerChannelID`
    )
);

CREATE TABLE `CustomerEvent` (
    `EventID` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `ClientEventID` bigint  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `EventID`
    )
);

CREATE TABLE `CustomerEventData` (
    `EventDataId` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `EventID` bigint UNSIGNED NOT NULL ,
    `ContentID` int UNSIGNED NOT NULL ,
    `CustomerID` bigint UNSIGNED NOT NULL ,
    `EventTypeID` smallint UNSIGNED NOT NULL ,
    `EventDate` timestamp  NOT NULL ,
    `Quantity` smallint UNSIGNED NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `EventDataId`
    )
);



CREATE TABLE `ContentPrice` (
    `ContentPriceID` mediumint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `ContentID` int UNSIGNED NOT NULL ,
    `Price` decimal(8,2) UNSIGNED NOT NULL ,
    `Currency` char(3)  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `ContentPriceID`
    )
);



ALTER TABLE `CustomerData` ADD CONSTRAINT `fk_CustomerData_CustomerID` FOREIGN KEY(`CustomerID`)
REFERENCES `Customer` (`CustomerID`);

ALTER TABLE `CustomerData` ADD CONSTRAINT `fk_CustomerData_ChannelTypeID` FOREIGN KEY(`ChannelTypeID`)
REFERENCES `ChannelType` (`ChannelTypeID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_EventID` FOREIGN KEY(`EventID`)
REFERENCES `CustomerEvent` (`EventID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_ContentID` FOREIGN KEY(`ContentID`)
REFERENCES `Content` (`ContentID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_CustomerID` FOREIGN KEY(`CustomerID`)
REFERENCES `Customer` (`CustomerID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_EventTypeID` FOREIGN KEY(`EventTypeID`)
REFERENCES `EventType` (`EventTypeID`);

ALTER TABLE `ContentPrice` ADD CONSTRAINT `fk_ContentPrice_ContentID` FOREIGN KEY(`ContentID`)
REFERENCES `Content` (`ContentID`);

