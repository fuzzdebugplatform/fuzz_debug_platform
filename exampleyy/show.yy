{
 function field(num)
     print(string.format("field%d", num))
 end
}

query:
    select_or_explain_select

select_or_explain_select:
    select

select:
   { num = 0 } SELECT DISTINCT select_list FROM _table where
   | { num = 0 } SELECT select_list FROM _table where

select_list:
   select_item AS { num=num+1; field(num) } |
   select_item AS { num=num+1; field(num) } , select_list

select_item:
   func | aggregate_func

aggregate_func:
   COUNT( func )
   | AVG( func )
   | SUM( func )
   | MAX( func )
   | MIN( func )
   | GROUP_CONCAT( func, func )
   | BIT_AND( arg )
    | BIT_COUNT( arg )
	| BIT_LENGTH( arg )
   | BIT_OR( arg )
   | BIT_XOR( arg )

where:
   | WHERE func ;

having:
   | HAVING func ;

order_by:
   | ORDER BY func | ORDER BY func, func ;

func:
    math_func |
	str_func |
	cast_oper

math_func:
   ABS( arg ) |
   PI( ) | POW( arg, arg ) | POWER( arg, arg ) |
   TAN( arg ) | TRUNCATE( arg, truncate_second_arg ) ;

str_func:
   ASCII( arg ) |
   BIN( arg ) |
   BIT_LENGTH( arg ) |
   UCASE( arg ) |
   UNHEX( arg ) |
   UPPER( arg )

truncate_second_arg:
   _digit | _digit | _tinyint_unsigned | arg ;

arg:
   _field | value | ( func ) ;

value:
   _int | _bigint | _smallint | _int_usigned | _letter | _english | _datetime | _date | _time | NULL

cast_oper:
   BINARY arg | CAST( arg AS type ) | CONVERT( arg, type ) | CONVERT( arg USING charset )

type:
   BINARY | BINARY(_digit)
   | CHAR | CHAR(_digit)
   | DATE
   | DATETIME | DECIMAL
   | DECIMAL(decimal_m)
   | DECIMAL(decimal_m,decimal_d) | SIGNED | TIME | UNSIGNED

charset:
   utf8 | latin1 ;

decimal_m:
    { decimal_m = math.random(1,65) }

decimal_d:
    { decimal_d = math.random(0, decimal_m) }