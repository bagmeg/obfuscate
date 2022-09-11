package obfuscate

import (
	"github.com/bagmeg/obfuscate/config"
	"testing"
)

func TestObfuscate(t *testing.T) {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		t.Errorf("Failed to load %v", "./config.yaml")
	}

	SQLobfuscator, _ := NewObfuscator(&cfg)

	if SQLobfuscator == nil {
		t.Error("Received nil SQLObfuscator")
	}

	var tests = []struct {
		query string
		want  string
	}{
		{
			"select * from user",
			"select * from user",
		},
		{
			"select * from user where age > 7",
			"select * from user where age > ?",
		},
		{
			"select age as AGE from user where age > 7",
			"select age as AGE from user where age > ?",
		},
		{
			"select age as AGE, number as NUM from user where age > 7",
			"select age as AGE, number as NUM from user where age > ?",
		},
		{
			"select age AS AGE from user",
			"select age AS AGE from user",
		},
		{
			`select age as 'AG''E' from user`,
			`select age as ? from user`,
		},
		{
			`select age as '\'user\'' from user`,
			`select age as ? from user`,
		},
		{
			`select age as 'user inside \'at\' home' from user`,
			`select age as ? from user`,
		},
		{
			`SELECT articles.* FROM articles WHERE articles.id = 1 LIMIT 1, 20`,
			"SELECT articles.* FROM articles WHERE articles.id = ? LIMIT ?, ?",
		},
		{
			"INSERT INTO `testTable` ( NAME ) VALUES (?)",
			"INSERT INTO testTable ( NAME ) VALUES ?",
		},
		{
			`SELECT * FROM orders WHERE customer_id = ? AND quantity > ?`,
			`SELECT * FROM orders WHERE customer_id = ? AND quantity > ?`,
		},
		{
			`DELETE FROM public.container_info WHERE id LIKE $1 ESCAPE $3 AND tenant_id LIKE $2 ESCAPE $4`,
			`DELETE FROM public.container_info WHERE id LIKE ? ESCAPE ? AND tenant_id LIKE ? ESCAPE ?`,
		},
		{
			`DELETE FROM public.container_info WHERE id LIKE $1 ESCAPE $tail AND tenant_id LIKE $2 ESCAPE $4`,
			`DELETE FROM public.container_info WHERE id LIKE ? ESCAPE $tail AND tenant_id LIKE ? ESCAPE ?`,
		},
		{
			"CREATE FUNCTION add(integer, integer) RETURNS integer\n AS 'select $1 + $2;'\n LANGUAGE SQL\n IMMUTABLE\n RETURNS NULL ON NULL INPUT;",
			"CREATE FUNCTION add ( integer, integer ) RETURNS integer AS ? LANGUAGE SQL IMMUTABLE RETURNS NULL ON NULL INPUT ;",
		},
		{
			"CREATE OR REPLACE FUNCTION public.retrieve_explain_plan_json(in_query text, OUT out_explain json)\r\n RETURNS SETOF json\r\n LANGUAGE plpgsql\r\n STRICT SECURITY DEFINER\r\nAS $function$\r\nBEGIN\r\n   RETURN QUERY EXECUTE 'EXPLAIN (FORMAT JSON) ' || in_query;\r\nEND;\r\n$function$",
			"CREATE OR REPLACE FUNCTION public.retrieve_explain_plan_json ( in_query text, OUT out_explain json ) RETURNS SETOF json LANGUAGE plpgsql STRICT SECURITY DEFINER AS $function$ BEGIN RETURN QUERY EXECUTE ? || in_query ; END ; $function$",
		},
		{
			`CREATE EVENT IF NOT EXISTS getDigest
ON SCHEDULE EVERY 1 SECOND
ENDS CURRENT_TIMESTAMP + INTERVAL 5 MINUTE ON COMPLETION NOT PRESERVE
DO
INSERT IGNORE INTO percona.digest_seen SELECT CURRENT_SCHEMA, DIGEST, SQL_TEXT FROM performance_schema.events_statements_history WHERE DIGEST IS NOT NULL GROUP BY current_schema, digest LIMIT 50;`,
			`CREATE EVENT IF NOT EXISTS getDigest ON SCHEDULE EVERY ? SECOND ENDS CURRENT_TIMESTAMP + INTERVAL ? MINUTE ON COMPLETION NOT PRESERVE DO INSERT IGNORE INTO percona.digest_seen SELECT CURRENT_SCHEMA, DIGEST, SQL_TEXT FROM performance_schema.events_statements_history WHERE DIGEST IS NOT NULL GROUP BY current_schema, digest LIMIT ? ;`,
		},
		{
			`SELECT max_session as max_connection
		 , total_session as total_session
		 , lock_session as lock_session
		 , total_session / max_session * 100.0 as connection_usage
		FROM (select current_setting('max_connections')::float as max_session
		       , (select count(*)::float from pg_stat_activity) as total_session
		       , (select count(*) from pg_stat_activity where cardinality(pg_blocking_pids(pid)) > 0) as lock_session
		) session`,
			"SELECT max_session as max_connection, total_session as total_session, lock_session as lock_session, total_session / max_session * ? as connection_usage FROM ( select current_setting ( ? ) :: float as max_session, ( select count ( * ) :: float from pg_stat_activity ) as total_session, ( select count ( * ) from pg_stat_activity where cardinality ( pg_blocking_pids ( pid ) ) > ? ) as lock_session ) session",
		},
	}
	for _, tt := range tests {
		t.Run("obfuscate_test", func(t *testing.T) {
			res, _ := SQLobfuscator.Scan(tt.query)
			if res != tt.want {
				t.Errorf("\ngot : %s\nwant: %s", res, tt.want)
			}
		})
	}
}
