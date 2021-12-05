grammar sed;

statements   
  : statement ( ';' statement )* ;

statement    
  : ':' label
  | address?  function | '{' functions '}' 
;

functions  
  : function ( ';'? function )* ;

address 
  : firstAddress ( ',' secondAddress )? '!'? ;

firstAddress
  : absoluteAddress ;

absoluteAddress
  : [1-9][0-9]*
  | '$'
  | pattern
;

secondAddress
  : absoluteAddress |  relativeAddress ;

relativeAddress
  : '+' [1-9][0-9]* ;

pattern        
  : '/' .+ '/' 
  | '\' (delimiter) .+ \1
;

delimiter
  : [^\\n] ;

function       
  : 'b' label?
  | 't' label?
  | 's' (delimiter) .* \1 .* \1 flags*
  | [nNgGhHxlpD]
;

label            
  : [a-z]+ ；

flags
  : [[0-9]gp]
