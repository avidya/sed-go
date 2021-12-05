A sed interpreter written in Golang.

Most behaviors of this program is implemented follow GNU-sed, except all those operations related to external files are kepted unfinished. 

the grammar is even more loosen compare to the GNU one: the semicolon `;` after each commands is not required(except for `s`). so we can write commands in this style:
```Bash
$ echo -e 'a\nb\nc' | ./sed -n 'hgNhnxGs/\n//g;s/abc/ABC/p'
ABC
```

Only two command line parameters are supplied.
```Bash
$ ./sed -h
Usage of ./sed:
  -n	inhibit printing
  -s	work in pseudo-stream mode
```

`-r` is not necessary since ERE is enabled as default, and there's no way to switch it off. As the standard Golang libs use RE2 as the underlying regular exression support, there may exists some subtle differences between these two standards. 

`-s` can be regarded as an extension upon the original version. 

basically spreaking, a sed implementation follows the POSIX standard is a line based editing tool. Yes, no matter what the name it is from. It's line based, not stream oriented! one will get into trouble when facing problem to process spanning lines, and that's why N, n, b, t these commands come into play.

think about this problem, if we want to remove all the newline feeds between a `<p>`, `</p>` pair. How to do that? 

```Bash
$ echo -e "a\n<p>\nb\nc\n</p>\nd\n<p>\ne\nf\n</p>\ng"
a
<p>
b
c
</p>
d
<p>
e
f
</p>
g
```

although there's multiple ways to achieve that, but I believe the solution listed below shows us an idea how a typical stream oriented processing way is like. 

```Bash
$ echo -e "a\n<p>\nb\nc\n</p>\nd\n<p>\ne\nf\n</p>\ng" | sed '/<p>/!b;:x N;/<\/p>/!bx;s/\n//g;'
a
<p>bc</p>
d
<p>ef</p>
g
```

 - step 1. before we find the start pattern, we simplly output everything that we read from stream.
 - step 2. after we find the start pattern, we are in a match range. all the inputs should be buffered until the stop pattern is founded. 
 - step 3. actually processing

Using this tool, you can achieve this just like the following:
```Bash
$ echo -e "a\n<p>\nb\nc\n</p>\nd\n<p>\ne\nf\n</p>\ng" | ./sed -s '/<p>/,/<\/p>/s/\n//g'
a
<p>bc</p>
d
<p>ef</p>
g
```

the stream mode is also applied to line pattern, and relative range pattern as well

remove all the `\n':
```Bash
$ echo -e "a\n<p>\nb\nc\n</p>\nd\n<p>\ne\nf\n</p>\ng" | ./sed -s '1,$s/\n//g'
a<p>bc</p>d<p>ef</p>g
```

you can even omit 1, $. when no address pattern is defined. this tool will take the first and  the last line as the range pattern.


```Bash
$ echo -e "a\n<p>\nb\nc\n</p>\nd\n<p>\ne\nf\n</p>\ng" | ./sed -s 's/\n//g'
a<p>bc</p>d<p>ef</p>g
```

Enjoy! any feedback: benimaur@gmail.com
