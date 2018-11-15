create user test with password 'test';
create database test owner test;
grant all privileges on database test to test;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO test;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO test;

# mini_user表
id:主键，自增
openid：微信用户的微信id标识
userinfo：用户基本信息
token：标识用户登录的token

> create table mini_user(id serial primary key,openid text,userinfo text,token text);


# mini_orders表
id：主键、自增
askuseropenid:提问者，用openid来标识用户
questiontitle:问题标题
questioncontent：问题的内容
questionreward:问题的红包
state:0:表示红包未付款成功 1：红包已经付款成功
out_trade_no:商户的唯一标识id
createtime:创建的时间戳

> create table mini_orders(id serial primary key,askuseropenid text,questiontitle text,questioncontent text,questionreward int,state int,createtime bigint,out_trade_no text);

# mini_answers 表
id：主键、自增
answeruseropenid:回答的问题的用户openid
answercontent:回答的内容
questionid：问题的id，对应的就是mini_orders表中的id
questiontitle:问题的标题，这个是方便我们查询数据而设置的。
createtime：创建的时间

> create table mini_answers(id serial primary
