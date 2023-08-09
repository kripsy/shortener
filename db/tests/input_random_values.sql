do $$
DECLARE v_original_url TEXT;
DECLARE v_short_url TEXT;
DECLARE v_user_id INT;

begin

   for inc in 1..10000 loop
       insert into users values (inc);
   end loop;


	for inc in 1..800000 loop
		v_original_url:= (SELECT md5(random()::text));
		v_short_url:= (SELECT md5(random()::text));
		v_user_id:=(select floor(random()* (500-1 + 1) + 1));

		insert into urls (id, original_url, short_url, user_id, is_deleted) 
		values(
			inc,
			v_original_url,
			v_short_url,
			v_user_id,random() > 0.5
			);
	end loop;
	
end; $$
