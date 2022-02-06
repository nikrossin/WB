create table orders(

                       order_uid varchar(50) primary key,
                       track_number varchar(50),
                       entry varchar(20),
                       delivery integer,
                       payment varchar(50),
                       items integer,
                       locale varchar(10),
                       internal_signature varchar(50),
                       customer_id varchar(20),
                       delivery_service varchar(20),
                       shardkey varchar(15),
                       sm_id integer,
                       date_created timestamp,
                       oof_shard varchar(15)
);

create table deliveries(
                           id serial primary key,
                           name varchar(30),
                           phone varchar(20),
                           zip varchar(20),
                           city varchar(20),
                           address varchar(25),
                           region varchar(20),
                           email varchar(30)
);

create table payments(
                         transaction varchar(50) primary key,
                         request_id varchar(50),
                         currency varchar(10),
                         provider varchar(20),
                         amount integer,
                         payment_dt bigint,
                         bank varchar(30),
                         delivery_cost integer,
                         goods_total integer,
                         custom_fee integer

);

create table itemstab(
                         id serial primary key,
                         chrt_id integer,
                         track_number varchar(50),
                         price integer,
                         rid varchar(50),
                         name varchar(20),
                         sale integer,
                         size varchar(10),
                         total_price integer,
                         nm_id integer,
                         brand varchar(20),
                         status integer
);


ALTER TABLE itemstab
    ADD order_id CHARACTER VARYING(50);

alter table orders
    add constraint fk_deliv
        foreign key (delivery) references deliveries(id);


alter table orders
    add constraint fk_pay
        foreign key (payment) references payments(transaction);

alter table itemstab
    add constraint fk_order
        foreign key (order_id) references orders(order_uid);




INSERT INTO deliveries(name,phone,zip,city,address,region,email)
VALUES ('Test','+79999','232323','moscow','izumrudnaya','mos','a@a');

INSERT INTO payments
VALUES ('saee9qq9w8q9w','212121','USD','wbpay',1500,129998000,'sber',1500,317,0);


INSERT INTO itemstab(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (99898989,'WBUHEUHIJE',435,'dsd8sdu8ud8d','Masrf',30,'0',317, 11323242,'zara',202);


CREATE FUNCTION addDatadelivery (name varchar(30),phone varchar(20),zip varchar(20),city varchar(20),address varchar(25),region varchar(20),email varchar(30)) returns void
as
    $$

	INSERT INTO deliveries(name,phone,zip,city,address,region,email)
	VALUES (name,phone,zip,city,address,region,email);
$$
language 'sql';


CREATE FUNCTION addDataPayments(transaction varchar(50),request_id varchar(50),currency varchar(10),provider varchar(20),amount integer,payment_dt bigint,bank varchar(30),delivery_cost integer,goods_total integer,custom_fee integer) returns void
as
    $$

	INSERT INTO payments(transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee)
	VALUES (transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee);
$$
language 'sql';



CREATE FUNCTION addDataItems(chrt_id integer,track_number varchar(50),price integer,rid varchar(50),name varchar(20),sale integer,size varchar(10),total_price integer,nm_id integer,brand varchar(20),status integer, order_id varchar(50)) returns void
as
    $$

	INSERT INTO itemstab(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id)
	VALUES (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id);
$$
language 'sql'

CREATE FUNCTION addDataOrders(order_uid varchar(50),track_number varchar(50),entry varchar(20),payment varchar(50),locale varchar(10),internal_signature varchar(50),customer_id varchar(20),delivery_service varchar(20),shardkey varchar(15),sm_id integer,date_created timestamp,oof_shard varchar(15)) returns void
as
    $$

	INSERT INTO orders(order_uid,track_number,entry,delivery,payment,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard)
	VALUES (order_uid,track_number,entry,
		(SELECT MAX(id) FROM deliveries),
		payment,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard);
$$
language 'sql'



