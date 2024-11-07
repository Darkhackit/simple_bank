CREATE TABLE "users" (
     "id" bigserial PRIMARY KEY ,
     "email" varchar UNIQUE NOT NULL ,
     "username" varchar UNIQUE NOT NULL ,
     "active" boolean  ,
     "name" varchar  ,
     "password" varchar NOT NULL ,
     "password_changed_at" timestamp DEFAULT (now()) ,
     "created_at" timestamp DEFAULT (now()) 
);
CREATE INDEX ON "users" ("id");

ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON account ("owner","currency");

ALTER TABLE "account" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner",currency)