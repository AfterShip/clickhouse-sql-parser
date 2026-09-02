# Benchmarks

```sh
go test -bench=. -benchmem ./parser
```

Results

```
$ go test -bench=. -benchmem ./parser
goos: darwin
goarch: arm64
pkg: github.com/AfterShip/clickhouse-sql-parser/parser
cpu: Apple M5
BenchmarkParseSQLFiles/access_tuple_with_dot.sql-10         	   92634	     13255 ns/op	   14125 B/op	     160 allocs/op
BenchmarkParseSQLFiles/create_window_view.sql-10            	  205540	      6027 ns/op	    7240 B/op	      50 allocs/op
BenchmarkParseSQLFiles/query_with_expr_compare.sql-10       	  142744	      7984 ns/op	    8024 B/op	      74 allocs/op
BenchmarkParseSQLFiles/select_case_multiple_when.sql-10     	  110898	     10199 ns/op	   10256 B/op	      36 allocs/op
BenchmarkParseSQLFiles/select_case_when_exists.sql-10       	   74604	     15759 ns/op	   19392 B/op	      49 allocs/op
BenchmarkParseSQLFiles/select_case_when_regexp.sql-10       	  196922	      6724 ns/op	    6624 B/op	      30 allocs/op
BenchmarkParseSQLFiles/select_cast.sql-10                   	   53528	     21601 ns/op	   29048 B/op	      86 allocs/op
BenchmarkParseSQLFiles/select_clause_keyword_as_column.sql-10         	  360277	      4502 ns/op	    3712 B/op	      30 allocs/op
BenchmarkParseSQLFiles/select_clause_keyword_as_only_column.sql-10    	 1000000	      1105 ns/op	     664 B/op	       8 allocs/op
BenchmarkParseSQLFiles/select_column_alias_string.sql-10              	  548592	      2232 ns/op	    1952 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_concat_expr.sql-10                      	  340119	      3556 ns/op	    3224 B/op	      29 allocs/op
BenchmarkParseSQLFiles/select_dynamic_subcolumn_type_hint.sql-10      	  646071	      2253 ns/op	    1888 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_end_as_column_name.sql-10               	   58902	     18615 ns/op	   19760 B/op	      55 allocs/op
BenchmarkParseSQLFiles/select_expr.sql-10                             	 1304924	       924.9 ns/op	     776 B/op	      10 allocs/op
BenchmarkParseSQLFiles/select_expr_keyword_as_column.sql-10           	  366388	      3818 ns/op	    3040 B/op	      27 allocs/op
BenchmarkParseSQLFiles/select_extract_with_regex.sql-10               	   44274	     28570 ns/op	   30976 B/op	     115 allocs/op
BenchmarkParseSQLFiles/select_function_keyword_args.sql-10            	   14722	     83592 ns/op	  111456 B/op	     421 allocs/op
BenchmarkParseSQLFiles/select_intersect.sql-10                        	  119796	     10306 ns/op	    8800 B/op	      68 allocs/op
BenchmarkParseSQLFiles/select_interval_as_column_name.sql-10          	   13432	     88537 ns/op	  101795 B/op	     404 allocs/op
BenchmarkParseSQLFiles/select_item_with_modifiers.sql-10              	  154108	      7649 ns/op	    6384 B/op	      80 allocs/op
BenchmarkParseSQLFiles/select_json_type.sql-10                        	   57916	     19456 ns/op	   27712 B/op	     101 allocs/op
BenchmarkParseSQLFiles/select_keyword_alias_no_as.sql-10              	  807847	      1530 ns/op	    1328 B/op	      15 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_alias.sql-10                 	  629100	      2397 ns/op	    1848 B/op	      18 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_column.sql-10                	  367125	      3266 ns/op	    2880 B/op	      28 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_last_column.sql-10           	  712323	      1726 ns/op	    1528 B/op	      18 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_only_column.sql-10           	 1000000	      1186 ns/op	     664 B/op	       8 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_only_column_semicolon.sql-10 	  572713	      1981 ns/op	    2712 B/op	      16 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_arithmetic.sql-10       	  158814	      7191 ns/op	    9872 B/op	      43 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_as_function_argument.sql-10         	   97154	     12409 ns/op	   19000 B/op	      72 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_before_operator_keyword.sql-10      	   94813	     15085 ns/op	   14672 B/op	      51 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_cast_and_subscript.sql-10           	  160732	      7412 ns/op	    8896 B/op	      59 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_in_where.sql-10                     	  631492	      2001 ns/op	    2200 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_multiple_items.sql-10               	  402571	      3387 ns/op	    3952 B/op	      28 allocs/op
BenchmarkParseSQLFiles/select_order_by_timestamp.sql-10                           	  359602	      3145 ns/op	    3728 B/op	      25 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_basic.sql-10                     	  180240	      6580 ns/op	    5768 B/op	      54 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_from_to.sql-10                   	  188487	      6880 ns/op	    6312 B/op	      57 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_interpolate.sql-10               	  130939	      7718 ns/op	    7608 B/op	      70 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_interpolate_no_columns.sql-10    	  202425	      6296 ns/op	    5992 B/op	      57 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_staleness.sql-10                 	  143965	     10189 ns/op	    8440 B/op	      55 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_step.sql-10                      	   73323	     17200 ns/op	   17152 B/op	      78 allocs/op
BenchmarkParseSQLFiles/select_regexp.sql-10                                       	  103813	     10507 ns/op	    5496 B/op	      35 allocs/op
BenchmarkParseSQLFiles/select_reserved_keyword_qualifier.sql-10                   	  352544	      4520 ns/op	    1888 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_signed_number_after_bracket.sql-10                  	  133626	     11372 ns/op	    5800 B/op	      59 allocs/op
BenchmarkParseSQLFiles/select_simple.sql-10                                       	   68014	     21406 ns/op	   11112 B/op	     114 allocs/op
BenchmarkParseSQLFiles/select_simple_field_alias.sql-10                           	  539674	      2567 ns/op	    1856 B/op	      21 allocs/op
BenchmarkParseSQLFiles/select_simple_with_bracket.sql-10                          	  238954	      4870 ns/op	    4048 B/op	      55 allocs/op
BenchmarkParseSQLFiles/select_simple_with_cte_with_column_aliases.sql-10          	  151998	      8068 ns/op	    9040 B/op	      56 allocs/op
BenchmarkParseSQLFiles/select_simple_with_group_by_with_cube_totals.sql-10        	  329010	      3557 ns/op	    3024 B/op	      33 allocs/op
BenchmarkParseSQLFiles/select_simple_with_is_not_null.sql-10                      	  161505	      7355 ns/op	    5776 B/op	      57 allocs/op
BenchmarkParseSQLFiles/select_simple_with_is_null.sql-10                          	  188919	      6307 ns/op	    4960 B/op	      54 allocs/op
BenchmarkParseSQLFiles/select_simple_with_limit.sql-10                            	  373185	      3455 ns/op	    2872 B/op	      24 allocs/op
BenchmarkParseSQLFiles/select_simple_with_top_clause.sql-10                       	  810536	      1488 ns/op	    1392 B/op	      15 allocs/op
BenchmarkParseSQLFiles/select_simple_with_with_clause.sql-10                      	  188566	      5899 ns/op	    5720 B/op	      68 allocs/op
BenchmarkParseSQLFiles/select_table_alias_without_keyword.sql-10                  	  329026	      3700 ns/op	    3272 B/op	      45 allocs/op
BenchmarkParseSQLFiles/select_table_function_arg_exprs.sql-10                     	   39381	     29223 ns/op	   27112 B/op	     263 allocs/op
BenchmarkParseSQLFiles/select_table_function_with_query.sql-10                    	  269560	      5319 ns/op	    4488 B/op	      47 allocs/op
BenchmarkParseSQLFiles/select_trailing_comma_before_from.sql-10                   	  292267	      4534 ns/op	    3288 B/op	      26 allocs/op
BenchmarkParseSQLFiles/select_trailing_comma_before_from_keyword_table.sql-10     	  403844	      3103 ns/op	    3272 B/op	      27 allocs/op
BenchmarkParseSQLFiles/select_when_condition.sql-10                               	  170295	      6957 ns/op	    7344 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_window_comprehensive.sql-10                         	   16311	     89190 ns/op	   77024 B/op	     545 allocs/op
BenchmarkParseSQLFiles/select_window_cte.sql-10                                   	   49759	     23842 ns/op	   21488 B/op	     146 allocs/op
BenchmarkParseSQLFiles/select_window_keyword_name_in_parens.sql-10                	  312724	      3966 ns/op	    3592 B/op	      38 allocs/op
BenchmarkParseSQLFiles/select_window_named_in_parens.sql-10                       	  280594	      3872 ns/op	    3448 B/op	      39 allocs/op
BenchmarkParseSQLFiles/select_window_named_reference_extensions.sql-10            	  156788	      7735 ns/op	    7816 B/op	      68 allocs/op
BenchmarkParseSQLFiles/select_window_params.sql-10                                	   80479	     14558 ns/op	   15056 B/op	     118 allocs/op
BenchmarkParseSQLFiles/select_with_distinct.sql-10                                	  635102	      2333 ns/op	    1824 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_with_distinct_keyword.sql-10                        	  853233	      1301 ns/op	    1248 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_with_distinct_on_dotted_columns.sql-10              	  334320	      4006 ns/op	    3608 B/op	      44 allocs/op
BenchmarkParseSQLFiles/select_with_distinct_on_keyword.sql-10                     	  667465	      1664 ns/op	    1720 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_with_global_join_locality.sql-10                    	   21967	     55176 ns/op	   55600 B/op	     620 allocs/op
BenchmarkParseSQLFiles/select_with_group_by.sql-10                                	  155880	      7751 ns/op	    8256 B/op	      79 allocs/op
BenchmarkParseSQLFiles/select_with_join_only.sql-10                               	  559351	      2129 ns/op	    1720 B/op	      25 allocs/op
BenchmarkParseSQLFiles/select_with_keyword_in_group_by.sql-10                     	  212350	      5679 ns/op	    5400 B/op	      58 allocs/op
BenchmarkParseSQLFiles/select_with_keyword_placeholder.sql-10                     	  487777	      2540 ns/op	    2712 B/op	      26 allocs/op
BenchmarkParseSQLFiles/select_with_left_join.sql-10                               	  276420	      4348 ns/op	    4896 B/op	      43 allocs/op
BenchmarkParseSQLFiles/select_with_literal_table_name.sql-10                      	  772923	      1598 ns/op	    1792 B/op	      16 allocs/op
BenchmarkParseSQLFiles/select_with_multi_array_and_inner_join.sql-10              	  116157	     10381 ns/op	   10344 B/op	     129 allocs/op
BenchmarkParseSQLFiles/select_with_multi_array_join.sql-10                        	  164395	      7332 ns/op	    8968 B/op	      70 allocs/op
BenchmarkParseSQLFiles/select_with_multi_except.sql-10                            	  273000	      4440 ns/op	    4600 B/op	      50 allocs/op
BenchmarkParseSQLFiles/select_with_multi_join.sql-10                              	   69586	     17127 ns/op	   17504 B/op	     114 allocs/op
BenchmarkParseSQLFiles/select_with_multi_line_comment.sql-10                      	  977642	      1290 ns/op	    1648 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_with_multi_union.sql-10                             	  546528	      2374 ns/op	    2392 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_with_multi_union_distinct.sql-10                    	  538717	      2334 ns/op	    2520 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_with_number_field.sql-10                            	  426824	      2951 ns/op	    2920 B/op	      34 allocs/op
BenchmarkParseSQLFiles/select_with_parenthesized_union.sql-10                     	   45453	     27087 ns/op	   30944 B/op	     192 allocs/op
BenchmarkParseSQLFiles/select_with_placeholder.sql-10                             	  683054	      1772 ns/op	    1720 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_with_query_parameter.sql-10                         	  129188	      9796 ns/op	   10848 B/op	     100 allocs/op
BenchmarkParseSQLFiles/select_with_settings_additional_table_filters.sql-10       	   80194	     15749 ns/op	   15360 B/op	     152 allocs/op
BenchmarkParseSQLFiles/select_with_single_quote_table.sql-10                      	 1201622	      1080 ns/op	     976 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_with_string_expr.sql-10                             	  543085	      2087 ns/op	    2288 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_with_union_distinct.sql-10                          	  417825	      3073 ns/op	    3376 B/op	      26 allocs/op
BenchmarkParseSQLFiles/select_with_variable.sql-10                                	  585704	      2406 ns/op	    2144 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_with_window_function.sql-10                         	   46557	     23717 ns/op	   32880 B/op	     101 allocs/op
BenchmarkParseSQLFiles/select_without_from_where.sql-10                           	  410430	      2897 ns/op	    3152 B/op	      29 allocs/op
BenchmarkParseSQLFiles/set_simple.sql-10                                          	  609226	      2078 ns/op	    3768 B/op	      25 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_with_multi_join.sql-10         	   68352	     18177 ns/op	   17504 B/op	     114 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_with_window_function.sql-10    	   52179	     23419 ns/op	   32880 B/op	     101 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_simple_with_with_clause.sql-10 	  210837	      5592 ns/op	    5720 B/op	      68 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_with_left_join.sql-10          	  282788	      4524 ns/op	    4896 B/op	      43 allocs/op
BenchmarkParseComplexQueries/testdata/benchdata/posthog_huge_0.sql-10             	     475	   2508498 ns/op	 3132779 B/op	   14928 allocs/op
BenchmarkParseComplexQueries/testdata/benchdata/posthog_huge_1.sql-10             	     598	   2049905 ns/op	 2628580 B/op	   12872 allocs/op
PASS
ok  	github.com/AfterShip/clickhouse-sql-parser/parser	142.620s
```
