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
BenchmarkParseSQLFiles/access_tuple_with_dot.sql-10         	   41946	     32091 ns/op	   14157 B/op	     160 allocs/op
BenchmarkParseSQLFiles/create_window_view.sql-10            	  151291	      7624 ns/op	    7272 B/op	      50 allocs/op
BenchmarkParseSQLFiles/query_with_expr_compare.sql-10       	  135202	      8740 ns/op	    7160 B/op	      73 allocs/op
BenchmarkParseSQLFiles/select_case_multiple_when.sql-10     	  191094	      6216 ns/op	    4912 B/op	      34 allocs/op
BenchmarkParseSQLFiles/select_case_when_exists.sql-10       	  127548	     10734 ns/op	    5600 B/op	      44 allocs/op
BenchmarkParseSQLFiles/select_case_when_regexp.sql-10       	  287572	      3867 ns/op	    2816 B/op	      27 allocs/op
BenchmarkParseSQLFiles/select_cast.sql-10                   	   91058	     11445 ns/op	    8856 B/op	      66 allocs/op
BenchmarkParseSQLFiles/select_clause_keyword_as_column.sql-10         	  493710	      2411 ns/op	    2272 B/op	      26 allocs/op
BenchmarkParseSQLFiles/select_clause_keyword_as_only_column.sql-10    	 1479050	       868.7 ns/op	     696 B/op	       8 allocs/op
BenchmarkParseSQLFiles/select_column_alias_string.sql-10              	  462427	      2358 ns/op	    1984 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_concat_expr.sql-10                      	  336728	      5401 ns/op	    3256 B/op	      29 allocs/op
BenchmarkParseSQLFiles/select_dynamic_subcolumn_type_hint.sql-10      	  257604	      3905 ns/op	    1920 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_end_as_column_name.sql-10               	  148753	      8255 ns/op	    4304 B/op	      44 allocs/op
BenchmarkParseSQLFiles/select_expr.sql-10                             	 1000000	      1002 ns/op	     808 B/op	      10 allocs/op
BenchmarkParseSQLFiles/select_expr_keyword_as_column.sql-10           	  479106	      3358 ns/op	    2048 B/op	      24 allocs/op
BenchmarkParseSQLFiles/select_extract_with_regex.sql-10               	   71292	     17824 ns/op	   11936 B/op	     110 allocs/op
BenchmarkParseSQLFiles/select_function_keyword_args.sql-10            	   20080	     51650 ns/op	   38464 B/op	     406 allocs/op
BenchmarkParseSQLFiles/select_intersect.sql-10                        	  155968	      7772 ns/op	    6784 B/op	      67 allocs/op
BenchmarkParseSQLFiles/select_interval_as_column_name.sql-10          	   22311	     49430 ns/op	   35352 B/op	     352 allocs/op
BenchmarkParseSQLFiles/select_item_with_modifiers.sql-10              	  165195	      6981 ns/op	    6416 B/op	      80 allocs/op
BenchmarkParseSQLFiles/select_json_type.sql-10                        	  112285	     11913 ns/op	   10272 B/op	      88 allocs/op
BenchmarkParseSQLFiles/select_keyword_alias_no_as.sql-10              	  788174	      1348 ns/op	    1104 B/op	      14 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_alias.sql-10                 	  585783	      1827 ns/op	    1880 B/op	      18 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_column.sql-10                	  403886	      2995 ns/op	    2912 B/op	      28 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_last_column.sql-10           	  860019	      1448 ns/op	    1304 B/op	      17 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_only_column.sql-10           	 1491062	       848.1 ns/op	     696 B/op	       8 allocs/op
BenchmarkParseSQLFiles/select_keyword_as_only_column_semicolon.sql-10 	 1290060	       880.4 ns/op	     696 B/op	       8 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_arithmetic.sql-10       	  304675	      3783 ns/op	    2992 B/op	      31 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_as_function_argument.sql-10         	  229147	      5205 ns/op	    3608 B/op	      47 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_before_operator_keyword.sql-10      	  238518	      5192 ns/op	    3504 B/op	      35 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_cast_and_subscript.sql-10           	  261912	      5058 ns/op	    3488 B/op	      47 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_in_where.sql-10                     	  575012	      1749 ns/op	    1464 B/op	      17 allocs/op
BenchmarkParseSQLFiles/select_keyword_operand_multiple_items.sql-10               	  623074	      2015 ns/op	    1680 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_order_by_timestamp.sql-10                           	  634240	      1907 ns/op	    1712 B/op	      17 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_basic.sql-10                     	  191694	      5874 ns/op	    5160 B/op	      52 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_from_to.sql-10                   	  194827	      6163 ns/op	    5704 B/op	      55 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_interpolate.sql-10               	  148472	      8786 ns/op	    7384 B/op	      69 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_interpolate_no_columns.sql-10    	  170904	      6282 ns/op	    6024 B/op	      57 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_staleness.sql-10                 	  238216	      5637 ns/op	    4376 B/op	      44 allocs/op
BenchmarkParseSQLFiles/select_order_by_with_fill_step.sql-10                      	  166125	      7995 ns/op	    5472 B/op	      51 allocs/op
BenchmarkParseSQLFiles/select_regexp.sql-10                                       	  447778	      2446 ns/op	    1688 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_reserved_keyword_qualifier.sql-10                   	  779988	      1700 ns/op	    1920 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_signed_number_after_bracket.sql-10                  	  173463	      6093 ns/op	    5640 B/op	      59 allocs/op
BenchmarkParseSQLFiles/select_simple.sql-10                                       	  101438	     12097 ns/op	   10504 B/op	     113 allocs/op
BenchmarkParseSQLFiles/select_simple_field_alias.sql-10                           	  611782	      2030 ns/op	    1888 B/op	      21 allocs/op
BenchmarkParseSQLFiles/select_simple_with_bracket.sql-10                          	  276570	      4635 ns/op	    4080 B/op	      55 allocs/op
BenchmarkParseSQLFiles/select_simple_with_cte_with_column_aliases.sql-10          	  193498	      6424 ns/op	    4976 B/op	      54 allocs/op
BenchmarkParseSQLFiles/select_simple_with_group_by_with_cube_totals.sql-10        	  271335	      4109 ns/op	    3056 B/op	      33 allocs/op
BenchmarkParseSQLFiles/select_simple_with_is_not_null.sql-10                      	  163068	      8633 ns/op	    5616 B/op	      57 allocs/op
BenchmarkParseSQLFiles/select_simple_with_is_null.sql-10                          	  124342	      8121 ns/op	    4672 B/op	      53 allocs/op
BenchmarkParseSQLFiles/select_simple_with_limit.sql-10                            	  279214	      4178 ns/op	    2904 B/op	      24 allocs/op
BenchmarkParseSQLFiles/select_simple_with_top_clause.sql-10                       	  678796	      1555 ns/op	    1424 B/op	      15 allocs/op
BenchmarkParseSQLFiles/select_simple_with_with_clause.sql-10                      	  195356	      6100 ns/op	    5496 B/op	      67 allocs/op
BenchmarkParseSQLFiles/select_table_alias_without_keyword.sql-10                  	  343621	      3535 ns/op	    3048 B/op	      44 allocs/op
BenchmarkParseSQLFiles/select_table_function_arg_exprs.sql-10                     	   44941	     27705 ns/op	   26056 B/op	     263 allocs/op
BenchmarkParseSQLFiles/select_table_function_with_query.sql-10                    	  247592	      4908 ns/op	    4264 B/op	      46 allocs/op
BenchmarkParseSQLFiles/select_trailing_comma_before_from.sql-10                   	  672033	      1888 ns/op	    1528 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_trailing_comma_before_from_keyword_table.sql-10     	  545592	      2195 ns/op	    1512 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_when_condition.sql-10                               	  343081	      3903 ns/op	    1616 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_window_comprehensive.sql-10                         	   13759	     73222 ns/op	   64640 B/op	     542 allocs/op
BenchmarkParseSQLFiles/select_window_cte.sql-10                                   	   66794	     17567 ns/op	   16336 B/op	     142 allocs/op
BenchmarkParseSQLFiles/select_window_keyword_name_in_parens.sql-10                	  305048	      4016 ns/op	    3624 B/op	      38 allocs/op
BenchmarkParseSQLFiles/select_window_named_in_parens.sql-10                       	  278389	      3977 ns/op	    3480 B/op	      39 allocs/op
BenchmarkParseSQLFiles/select_window_named_reference_extensions.sql-10            	  125362	     10458 ns/op	    7848 B/op	      68 allocs/op
BenchmarkParseSQLFiles/select_window_params.sql-10                                	   78658	     14711 ns/op	   15088 B/op	     118 allocs/op
BenchmarkParseSQLFiles/select_with_distinct.sql-10                                	  552853	      2407 ns/op	    1856 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_with_distinct_keyword.sql-10                        	  905398	      1306 ns/op	    1280 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_with_distinct_on_dotted_columns.sql-10              	  339792	      4032 ns/op	    3576 B/op	      44 allocs/op
BenchmarkParseSQLFiles/select_with_distinct_on_keyword.sql-10                     	  549628	      1900 ns/op	    1752 B/op	      20 allocs/op
BenchmarkParseSQLFiles/select_with_global_join_locality.sql-10                    	   19771	     58056 ns/op	   53776 B/op	     619 allocs/op
BenchmarkParseSQLFiles/select_with_group_by.sql-10                                	  137679	      8616 ns/op	    8288 B/op	      79 allocs/op
BenchmarkParseSQLFiles/select_with_join_only.sql-10                               	  602020	      2157 ns/op	    1752 B/op	      25 allocs/op
BenchmarkParseSQLFiles/select_with_keyword_in_group_by.sql-10                     	  193964	      7432 ns/op	    5432 B/op	      58 allocs/op
BenchmarkParseSQLFiles/select_with_keyword_placeholder.sql-10                     	  385830	      2951 ns/op	    2488 B/op	      25 allocs/op
BenchmarkParseSQLFiles/select_with_left_join.sql-10                               	  249094	      5832 ns/op	    4928 B/op	      43 allocs/op
BenchmarkParseSQLFiles/select_with_literal_table_name.sql-10                      	  644481	      2233 ns/op	    1824 B/op	      16 allocs/op
BenchmarkParseSQLFiles/select_with_multi_array_and_inner_join.sql-10              	  102230	     11334 ns/op	   10184 B/op	     130 allocs/op
BenchmarkParseSQLFiles/select_with_multi_array_join.sql-10                        	  163647	      6441 ns/op	    5288 B/op	      66 allocs/op
BenchmarkParseSQLFiles/select_with_multi_except.sql-10                            	  235374	      5238 ns/op	    4632 B/op	      50 allocs/op
BenchmarkParseSQLFiles/select_with_multi_join.sql-10                              	  115492	     10876 ns/op	    9216 B/op	      94 allocs/op
BenchmarkParseSQLFiles/select_with_multi_line_comment.sql-10                      	  839001	      1636 ns/op	    1680 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_with_multi_union.sql-10                             	  485973	      2470 ns/op	    2424 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_with_multi_union_distinct.sql-10                    	  506588	      2576 ns/op	    2552 B/op	      19 allocs/op
BenchmarkParseSQLFiles/select_with_number_field.sql-10                            	  380826	      3483 ns/op	    2632 B/op	      33 allocs/op
BenchmarkParseSQLFiles/select_with_parenthesized_union.sql-10                     	   41083	     28368 ns/op	   28928 B/op	     191 allocs/op
BenchmarkParseSQLFiles/select_with_placeholder.sql-10                             	  701680	      1788 ns/op	    1496 B/op	      18 allocs/op
BenchmarkParseSQLFiles/select_with_query_parameter.sql-10                         	  115622	     10473 ns/op	   10368 B/op	      99 allocs/op
BenchmarkParseSQLFiles/select_with_settings_additional_table_filters.sql-10       	   69403	     18035 ns/op	   15392 B/op	     152 allocs/op
BenchmarkParseSQLFiles/select_with_single_quote_table.sql-10                      	 1000000	      1383 ns/op	    1008 B/op	      13 allocs/op
BenchmarkParseSQLFiles/select_with_string_expr.sql-10                             	  412788	      2446 ns/op	    2320 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_with_union_distinct.sql-10                          	  356832	      3263 ns/op	    3408 B/op	      26 allocs/op
BenchmarkParseSQLFiles/select_with_variable.sql-10                                	  563797	      2305 ns/op	    2176 B/op	      23 allocs/op
BenchmarkParseSQLFiles/select_with_window_function.sql-10                         	   91530	     13406 ns/op	   11280 B/op	      92 allocs/op
BenchmarkParseSQLFiles/select_without_from_where.sql-10                           	  391298	      3189 ns/op	    2928 B/op	      28 allocs/op
BenchmarkParseSQLFiles/set_simple.sql-10                                          	  495375	      2762 ns/op	    3800 B/op	      25 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_with_multi_join.sql-10         	  114853	     10731 ns/op	    9216 B/op	      94 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_with_window_function.sql-10    	   86808	     13809 ns/op	   11280 B/op	      92 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_simple_with_with_clause.sql-10 	  192193	      6218 ns/op	    5496 B/op	      67 allocs/op
BenchmarkParseComplexQueries/testdata/query/select_with_left_join.sql-10          	  246328	      4892 ns/op	    4928 B/op	      43 allocs/op
BenchmarkParseComplexQueries/testdata/benchdata/posthog_huge_0.sql-10             	     660	   1938482 ns/op	 1313672 B/op	   14554 allocs/op
BenchmarkParseComplexQueries/testdata/benchdata/posthog_huge_1.sql-10             	     759	   1579115 ns/op	 1106177 B/op	   12559 allocs/op
PASS
ok  	github.com/AfterShip/clickhouse-sql-parser/parser	139.761s
```
