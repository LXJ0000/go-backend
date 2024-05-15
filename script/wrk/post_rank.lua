wrk.method="GET"
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicm9vdCIsImlkIjoxODE0MjQ5MzExNDQyMDgzODQsImV4cCI6MTcxNTc5MzQ1N30.AI72bGlPyqus8u_buDevXZ-snIM5Vswl5j4CGAIc_s0"
-- wrk -t4 -d5s -c50 -s ./script/wrk/post_rank.lua http://localhost:8080/api/post/rank