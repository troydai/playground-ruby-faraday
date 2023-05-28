require 'faraday'
require 'faraday/net_http_persistent'

puts "faraday client"

class Client
    def initialize
    end

    def tell
        begin
            conn.get '/tell' 
        rescue Faraday::ConnectionFailed => e
            puts "Connection failed: #{e}"
        end
    end

    private

    def conn
        @conn ||= new_connect
    end

    def new_connect
        Faraday.new(:url => 'http://localhost:8080') do |f|
            f.adapter :net_http_persistent, pool_size: 5 do |http|
              http.idle_timeout = 100
            end
        end
    end
end

starting = Process.clock_gettime(Process::CLOCK_MONOTONIC)

c = Client.new
for i in 1..100
    c.tell
end

ending = Process.clock_gettime(Process::CLOCK_MONOTONIC)
elapsed = ending - starting
puts elapsed 