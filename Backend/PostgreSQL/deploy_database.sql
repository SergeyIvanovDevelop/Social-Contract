DROP SEQUENCE IF EXISTS "public"."alert_log_log_id_seq";
CREATE SEQUENCE "public"."alert_log_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

DROP SEQUENCE IF EXISTS "public"."cards_card_id_seq";
CREATE SEQUENCE "public"."cards_card_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

DROP SEQUENCE IF EXISTS "public"."contracts_contract_id_seq";
CREATE SEQUENCE "public"."contracts_contract_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

DROP SEQUENCE IF EXISTS "public"."users_user_id_seq";
CREATE SEQUENCE "public"."users_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

DROP TABLE IF EXISTS "public"."alert_log";
CREATE TABLE "public"."alert_log" (
  "log_id" int8 NOT NULL DEFAULT nextval('alert_log_log_id_seq'::regclass),
  "user_id" int8,
  "contract_id" int8,
  "alert_timestamp" timestamp(6)
)
;

DROP TABLE IF EXISTS "public"."cards";
CREATE TABLE "public"."cards" (
  "card_id" int8 NOT NULL DEFAULT nextval('cards_card_id_seq'::regclass),
  "user_id" int8,
  "Name" varchar COLLATE "pg_catalog"."default",
  "PersonalAcc" varchar COLLATE "pg_catalog"."default",
  "BankName" varchar COLLATE "pg_catalog"."default",
  "BIC" varchar COLLATE "pg_catalog"."default",
  "CorrespAcc" varchar COLLATE "pg_catalog"."default",
  "KPP" varchar COLLATE "pg_catalog"."default",
  "PayeeINN" varchar COLLATE "pg_catalog"."default"
)
;

DROP TABLE IF EXISTS "public"."contracts";
CREATE TABLE "public"."contracts" (
  "contract_id" int8 NOT NULL DEFAULT nextval('contracts_contract_id_seq'::regclass),
  "payer" int8,
  "recipient" int8,
  "contract_name" varchar(255) COLLATE "pg_catalog"."default",
  "checker_time" varchar(255),
  "completed" bool DEFAULT false
)
;

DROP TABLE IF EXISTS "public"."diary_credentials";
CREATE TABLE "public"."diary_credentials" (
  "user_id" int8 NOT NULL,
  "diary_login" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "diary_password" varchar(255) COLLATE "pg_catalog"."default"
)
;

DROP TABLE IF EXISTS "public"."prices";
CREATE TABLE "public"."prices" (
  "contract_id" int8,
  "grade" int2,
  "price" int2
)
;

DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "user_id" int8 NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
  "login" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "hash" varchar(255) COLLATE "pg_catalog"."default",
  "registration_timestamp" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "surname" varchar(255) COLLATE "pg_catalog"."default",
  "firstname" varchar(255) COLLATE "pg_catalog"."default",
  "photo" text
)
;

DROP VIEW IF EXISTS "public"."view_price_for_contracts";
CREATE VIEW "public"."view_price_for_contracts" AS  SELECT prices.contract_id,
    prise_1.price AS g1,
    prise_2.price AS g2,
    prise_3.price AS g3,
    prise_4.price AS g4,
    prise_5.price AS g5
   FROM prices
     LEFT JOIN prices prise_1 ON prise_1.grade = '1'::smallint
     LEFT JOIN prices prise_2 ON prise_2.grade = '2'::smallint
     LEFT JOIN prices prise_3 ON prise_3.grade = '3'::smallint
     LEFT JOIN prices prise_4 ON prise_4.grade = '4'::smallint
     LEFT JOIN prices prise_5 ON prise_5.grade = '5'::smallint
  GROUP BY prices.contract_id, prise_1.price, prise_2.price, prise_3.price, prise_4.price, prise_5.price;

DROP VIEW IF EXISTS "public"."view_auth_data";
CREATE VIEW "public"."view_auth_data" AS  SELECT users.login,
    users.hash
   FROM users;

DROP VIEW IF EXISTS "public"."view_contracts";
CREATE VIEW "public"."view_contracts" AS  SELECT payer.login AS payer_login,
    contracts.contract_name,
    recipient.surname AS recipient_surname,
    recipient.firstname AS recipient_firstname,
    recipient.photo AS recipient_photo
   FROM users payer
     RIGHT JOIN contracts ON payer.user_id = contracts.payer
     LEFT JOIN users recipient ON contracts.recipient = recipient.user_id;

DROP VIEW IF EXISTS "public"."view_cards";
CREATE VIEW "public"."view_cards" AS  SELECT users.login,
    cards."Name",
    cards."PersonalAcc",
    cards."BankName",
    cards."BIC",
    cards."CorrespAcc",
    cards."PayeeINN",
    cards."KPP"
   FROM cards
     LEFT JOIN users USING (user_id);

DROP VIEW IF EXISTS "public"."view_contracts_with_price";
CREATE VIEW "public"."view_contracts_with_price" AS  SELECT payer.login AS payer_login,
    contracts.contract_name,
    recipient.surname AS recipient_surname,
    recipient.firstname AS recipient_firstname,
    recipient.photo AS recipient_photo,
    price.g1,
    price.g2,
    price.g3,
    price.g4,
    price.g5
   FROM users payer
     RIGHT JOIN contracts ON payer.user_id = contracts.payer
     LEFT JOIN users recipient ON contracts.recipient = recipient.user_id
     LEFT JOIN ( SELECT prices.contract_id,
            prise_1.price AS g1,
            prise_2.price AS g2,
            prise_3.price AS g3,
            prise_4.price AS g4,
            prise_5.price AS g5
           FROM prices
             LEFT JOIN prices prise_1 ON prise_1.grade = '1'::smallint AND prices.contract_id = prise_1.contract_id
             LEFT JOIN prices prise_2 ON prise_2.grade = '2'::smallint AND prices.contract_id = prise_2.contract_id
             LEFT JOIN prices prise_3 ON prise_3.grade = '3'::smallint AND prices.contract_id = prise_3.contract_id
             LEFT JOIN prices prise_4 ON prise_4.grade = '4'::smallint AND prices.contract_id = prise_4.contract_id
             LEFT JOIN prices prise_5 ON prise_5.grade = '5'::smallint AND prices.contract_id = prise_5.contract_id
          GROUP BY prices.contract_id, prise_1.price, prise_2.price, prise_3.price, prise_4.price, prise_5.price) price ON contracts.contract_id = price.contract_id;

DROP VIEW IF EXISTS "public"."view_diary_credentials";
CREATE VIEW "public"."view_diary_credentials" AS  SELECT users.login,
    diary_credentials.diary_login,
    diary_credentials.diary_password
   FROM diary_credentials
     LEFT JOIN users USING (user_id);

DROP VIEW IF EXISTS "public"."view_recipients_for_payers";
CREATE VIEW "public"."view_recipients_for_payers" AS  SELECT payer.login AS payer_login,
    payer.surname AS payer_surname,
    payer.firstname AS payer_firstname,
    payer.photo AS payer_photo,  
    recipient.surname AS recipient_surname,
    recipient.firstname AS recipient_firstname,
    recipient.photo AS recipient_photo,
    cards."Name",
    cards."PersonalAcc",
    cards."BankName",
    cards."BIC",
    cards."CorrespAcc",
    cards."KPP",
    cards."PayeeINN"
   FROM contracts
     LEFT JOIN users payer ON contracts.payer = payer.user_id
     LEFT JOIN users recipient ON contracts.recipient = recipient.user_id
     LEFT JOIN cards ON recipient.user_id = cards.user_id
  GROUP BY payer.login, payer.surname, payer.firstname, payer.photo, recipient.surname, recipient.firstname, recipient.photo, cards."Name", cards."PersonalAcc", cards."BankName", cards."BIC", cards."CorrespAcc", cards."KPP", cards."PayeeINN";

DROP VIEW IF EXISTS "public"."view_last_alerts";
CREATE VIEW "public"."view_last_alerts" AS  SELECT users.login,
    max(alert_log.alert_timestamp) AS max
   FROM alert_log
     LEFT JOIN users USING (user_id)
  GROUP BY users.login;

DROP VIEW IF EXISTS "public"."view_prices";
CREATE VIEW "public"."view_prices" AS  SELECT users.login AS recipient_login,
    contracts.contract_id,
    prices.grade,
    prices.price
   FROM prices
     LEFT JOIN contracts USING (contract_id)
     LEFT JOIN users ON users.user_id = contracts.recipient;

ALTER SEQUENCE "public"."alert_log_log_id_seq"
OWNED BY "public"."alert_log"."log_id";
SELECT setval('"public"."alert_log_log_id_seq"', 4, true);

ALTER SEQUENCE "public"."cards_card_id_seq"
OWNED BY "public"."cards"."card_id";
SELECT setval('"public"."cards_card_id_seq"', 3, true);

ALTER SEQUENCE "public"."contracts_contract_id_seq"
OWNED BY "public"."contracts"."contract_id";
SELECT setval('"public"."contracts_contract_id_seq"', 6, true);

ALTER SEQUENCE "public"."users_user_id_seq"
OWNED BY "public"."users"."user_id";
SELECT setval('"public"."users_user_id_seq"', 5, true);

ALTER TABLE "public"."alert_log" ADD CONSTRAINT "alert_log_pkey" PRIMARY KEY ("log_id");

ALTER TABLE "public"."cards" ADD CONSTRAINT "cards_pkey" PRIMARY KEY ("card_id");

ALTER TABLE "public"."contracts" ADD CONSTRAINT "contracts_pkey" PRIMARY KEY ("contract_id");

ALTER TABLE "public"."diary_credentials" ADD CONSTRAINT "diary_credentials_diary_login_key" UNIQUE ("diary_login");

ALTER TABLE "public"."diary_credentials" ADD CONSTRAINT "diary_credentials_pkey" PRIMARY KEY ("user_id", "diary_login");

ALTER TABLE "public"."prices" ADD CONSTRAINT "prices_contract_id_grade_key" UNIQUE ("contract_id", "grade");

ALTER TABLE "public"."users" ADD CONSTRAINT "login" UNIQUE ("login");

ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("user_id");

ALTER TABLE "public"."alert_log" ADD CONSTRAINT "alert_log_contract_id_fkey" FOREIGN KEY ("contract_id") REFERENCES "public"."contracts" ("contract_id") ON DELETE NO ACTION ON UPDATE CASCADE;
ALTER TABLE "public"."alert_log" ADD CONSTRAINT "alert_log_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON DELETE NO ACTION ON UPDATE CASCADE;

ALTER TABLE "public"."cards" ADD CONSTRAINT "cards_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON DELETE NO ACTION ON UPDATE CASCADE;

ALTER TABLE "public"."contracts" ADD CONSTRAINT "contracts_payer_fkey" FOREIGN KEY ("payer") REFERENCES "public"."users" ("user_id") ON DELETE NO ACTION ON UPDATE CASCADE;
ALTER TABLE "public"."contracts" ADD CONSTRAINT "contracts_recipient_fkey" FOREIGN KEY ("recipient") REFERENCES "public"."users" ("user_id") ON DELETE NO ACTION ON UPDATE CASCADE;

ALTER TABLE "public"."diary_credentials" ADD CONSTRAINT "diary_credentials_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON DELETE NO ACTION ON UPDATE CASCADE;

ALTER TABLE "public"."prices" ADD CONSTRAINT "prices_contract_id_fkey" FOREIGN KEY ("contract_id") REFERENCES "public"."contracts" ("contract_id") ON DELETE NO ACTION ON UPDATE CASCADE;
