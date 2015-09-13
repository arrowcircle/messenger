namespace :db do
  task :migrate do
    `migrate -url #{driver} -path ./migrations up`
  end

  task :rollback do
    `migrate -url #{driver} -path ./migrations down`
  end

  task :g, :name do |t, args|
    `migrate -url #{driver} -path ./migrations create #{args[:name]}`
  end

  def driver
    "postgres://ionbuggy@localhost/chat_development?sslmode=disable"
  end
end
