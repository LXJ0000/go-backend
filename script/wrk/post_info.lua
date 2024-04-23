wrk.method="GET"
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicm9vdCIsImlkIjoxNzMyNTQ4NDUwMTEwNzA5NzYsImV4cCI6MTcxMzg0NTU3M30.5Le2K7BFusWNvE_vObz1nDAiHEb_wJIOXGzlqhmGhYU"
-- wrk -t4 -d5s -c50 -s ./script/wrk/post_info.lua http://localhost:8080/api/post?post_id=173255181931122688