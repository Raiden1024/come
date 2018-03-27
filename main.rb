#!/usr/bin/env ruby
require 'httparty'
require 'json'
require 'optparse'
require 'notify-send'

def create_conf
  begin
    unless File.exist?("#{ENV['HOME']}/.config/comecfg.json")
      comecfg = {}
      begin
        print 'Openfire Server Address ? (format: http(s)://hostname(:port): '
        server_name = gets.chomp
        HTTParty.get("#{server_name}/plugins/restapi/v1/sessions/")
      rescue SocketError, Errno::EHOSTUNREACH, Errno::ECONNREFUSED, URI::InvalidURIError
        puts "\e[41mIncorrect Openfire Server, retry...\e[0m"
        retry
      end
      print 'openfire server secret key ? '
      secret_key = gets.chomp
      comecfg['server'] = server_name
      comecfg['firetoken'] = secret_key
      File.open("#{ENV['HOME']}/.config/comecfg.json", "w") do |f|
        f.write(comecfg.to_json)
      end
    end
    cfg = File.read("#{ENV['HOME']}/.config/comecfg.json")
    comecfg = JSON.parse(cfg)
  rescue Interrupt
    abort "\n\e[45mCanceled by user\e[0m"
  end

end

def get_API(key, server)
  url = "#{server}/plugins/restapi/v1/sessions/"
  headers = {
    "Accept" => 'application/json',
    "Authorization" => key
  }
  resp = HTTParty.get(url, headers: headers)
  json = resp.body
  begin
    JSON.parse(json)
  rescue JSON::ParserError
    puts "Incorrect API Key in #{ENV['HOME']}/.config/come.cfg"
  end

end

def list_sessions(sessions)
  begin
    sessions.each do |k, v|
      sessions[k].each do |v|
        print "#{v['username']}  #{v['hostAddress']} \n"
      end
    end
  rescue NoMethodError
    exit
  end
end

def show_IP(sessions, user)
  begin
    user_ip = ""
    sessions.each do |k, v|
      sessions[k].each do |v|
        user_ip = v['hostAddress'] if v['username'] == user
      end
    end
    return user_ip
  rescue NoMethodError
    exit
  end
end

def ssh_connect(ip)
  exec("ssh root@#{ip}")
end

def sessions_menu(sessions)
  s_username = []
  s_ip = []
  menu_count = 0

  sessions.each do |k, v|
    sessions[k].each do |v|
      s_username[menu_count] = v['username']
      s_ip[menu_count] = v['hostAddress']
      menu_count += 1
    end
  end
  s_username.each_with_index do |v, i|
    printf "%-4s - %-25s", "#{i+1}", "#{s_username[i]}"
    if i %2 == 1
      print "\n"
    end
  end
  puts "\n\e[45mConnected users: #{s_username.length} \e[0m\n"
  begin
    print "choose user number for ssh connection: "
    user_number = Integer(gets.chomp) rescue false
    if user_number
      if s_ip[user_number]
        ssh_connect(s_ip[user_number-1])
      else
        puts 'Invalid choice: number not in list'
      end

    else
      puts 'Invalid choice: you must type a number'
    end
  rescue Interrupt
    abort "\n\e[45mCanceled by user\e[0m"
  end
end

token = create_conf['firetoken']
server  = create_conf['server']
ARGV << '-h' if ARGV.empty?

OptionParser.new do |opts|
  opts.banner = 'Usage: come [options]'
  opts.on('-l', '--list', 'Display active users sessions list'){
    if ARGV[0]
      puts "\e[41mToo many parameters!\e[0m"
      puts opts
    else
      list_sessions(get_API(token, server))
    end
  }
  opts.on('-m', '--sessions-menu', 'Display interactive users sessions menu'){
    if ARGV[0]
      puts "\e[41mToo many parameters!\e[0m"
      puts opts
    else
      sessions_menu(get_API(token, server))
    end
  }
  opts.on('-i', '--ip', 'Display User IP Address, ex: come -i <user>'){
    if !ARGV[0]
      puts "\e[41mUsername missing!\e[0m"
    elsif show_IP(get_API(token, server), ARGV[0]) == ''
      puts "\e[41mUnknown User or User not connected!\e[0m"
    else
      puts show_IP(get_API(token, server), ARGV[0])
    end
  }
  opts.on('-c', '--connect', 'SSH connect to a user machine, ex: come -c <user>'){
    if !ARGV[0]
      puts "\e[41mUsername missing!\e[0m"
    elsif show_IP(get_API(token, server), ARGV[0]) == ''
      puts "\e[41mUnknown User or User not connected!\e[0m"
    else
      ssh_connect(show_IP(get_API(token, server), ARGV[0]))
    end
  }

  opts.on('-w', '--waiting', 'Wait for user Online status, ex: come -w <user>'){
    if !ARGV[0]
      puts '\e[41mUsername missing!\e[0m'
    else
      puts "Waiting for user status...\nCtrl-c to cancel"

      loop do
        begin
          break if show_IP(get_API(token, server), ARGV[0]) != ''
        rescue Interrupt
          abort "\n\e[45mCanceled by user\e[0m"
        end
      end
      puts "#{ARGV[0]}: ONLINE"
      NotifySend.send "COME Notify", "#{ARGV[0]}: ONLINE"
      exec("paplay /usr/local/lib/come/glass.ogg")
    end
  }
  opts.on('-v', '--version', 'Print version'){ puts 'COME (COnnectME) - version: 1.2'}
  opts.on('-h','--help', 'This help') do
    puts opts
  end
end.parse!
