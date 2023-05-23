create table if not exists Currencies(
    ID   varchar(100) not null,
    NumCode  integer  ,
    CharCode varchar(30) ,
    Nom      integer   ,
    Name     varchar(90),
    Value    float,
    Date date,
    primary key (ID, Date)
    );
