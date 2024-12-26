SELECT tuple('a','b','c').3, .1234;

SELECT toTypeName( tuple('a' as first,'b' as second ,'c' as third)::Tuple(first String,second String,third String)),
       (tuple('a' as first,'b' as second ,'c' as third)::Tuple(first String,second String,third String)).second,
       tuple('a','b','c').3,
       tupleElement(tuple('a','b','c'),1)