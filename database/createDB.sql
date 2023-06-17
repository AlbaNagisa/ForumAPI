CREATE TABLE "Categorie_Message" (
	"message_id"	INTEGER NOT NULL,
	"categorie_id"	INTEGER NOT NULL,
	FOREIGN KEY("message_id") REFERENCES "Message"("id"),
	FOREIGN KEY("categorie_id") REFERENCES "Categories"("id")
);
CREATE TABLE "Categories" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
INSERT INTO "Categories" VALUES (1,'Anxiety');
INSERT INTO "Categories" VALUES (2,'Depression');
INSERT INTO "Categories" VALUES (3,'Stress');
CREATE TABLE "Citation_Message" (
	"message_id"	INTEGER NOT NULL,
	"citation_id"	INTEGER NOT NULL,
	FOREIGN KEY("message_id") REFERENCES "Message"("id"),
	FOREIGN KEY("citation_id") REFERENCES "Message"("id")
);
CREATE TABLE "Image" (
	"id"	INTEGER NOT NULL,
	"path"	BLOB NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE "Image_Message" (
	"message_id"	INTEGER NOT NULL,
	"image_id"	INTEGER NOT NULL,
	FOREIGN KEY("message_id") REFERENCES "Message"("id"),
	FOREIGN KEY("image_id") REFERENCES "Image"("id")
);

CREATE TABLE "Message" (
	"id"	INTEGER NOT NULL,
	"author_id"	BIGINT NOT NULL,
	"date"	DATETIME NOT NULL,
	"title" TEXT NOT NULL,
	"content"	TEXT NOT NULL,
	"is_response"	BOOLEAN NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("author_id") REFERENCES "User"("id")
);
CREATE TABLE "Message_Prompt" (
	"message_id"	INTEGER NOT NULL,
	"prompt_id"	INTEGER NOT NULL,
	FOREIGN KEY("message_id") REFERENCES "Message"("id"),
	FOREIGN KEY("prompt_id") REFERENCES "Prompt"("id")
);
CREATE TABLE "Message_Response" (
	"message_id"	INTEGER NOT NULL,
	"response_id"	INTEGER NOT NULL,
	FOREIGN KEY("message_id") REFERENCES "Message"("id"),
	FOREIGN KEY("response_id") REFERENCES "Message"("id")
);
CREATE TABLE "Prompt" (
	"id"	INTEGER NOT NULL,
	"prompt"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE "SavedPost" (
	"user_id"	INTEGER NOT NULL,
	"message_id"	INTEGER NOT NULL,
	FOREIGN KEY("user_id") REFERENCES "User"("id"),
	FOREIGN KEY("message_id") REFERENCES "Message"("id")
);
CREATE TABLE "User" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"profileImage"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("profileImage") REFERENCES "Image"("id")
);
CREATE TABLE "Vote" (
	"user_id"	INTEGER NOT NULL,
	"message_id"	INTEGER NOT NULL,
	"vote"	NUMERIC NOT NULL,
	FOREIGN KEY("user_id") REFERENCES "User"("id"),
	FOREIGN KEY("message_id") REFERENCES "Message"("id")
);


DELETE FROM User WHERE id = 16;