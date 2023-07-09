INSERT INTO public.accounts
(allow_negative, balance)
VALUES(false, 0.00);

select account_id, sum(amount)
from transfers t 
where "status" = 'APPROVED'
group by account_id 

select min(updated_at - created_at) as min_time, max(updated_at - created_at) as max_time, avg(updated_at - created_at) as avg_time
from transfers t 
where "status" = 'APPROVED'

SELECT gs.updated_at, count(*)
FROM (
  SELECT
    date_trunc('second',min(updated_at)) AS min,
    date_trunc('second',max(updated_at))+interval '1 second' AS max
  FROM transfers
) AS t
CROSS JOIN LATERAL generate_series(t.min,t.max,'2 second') AS gs(updated_at)
INNER JOIN transfers AS f
  ON (EXTRACT(epoch FROM f.updated_at)::int/1) = (EXTRACT(epoch FROM gs.updated_at)::int/1)
GROUP BY gs.updated_at
order by count desc
