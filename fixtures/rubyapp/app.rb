ENV['RACK_ENV'] ||= 'development'

require 'rubygems'
require 'sinatra/base'
require 'json'
require 'typhoeus'

Bundler.require :default, ENV['RACK_ENV'].to_sym

class AppServer < Sinatra::Base
  get '/' do
    "Hi, I'm an app with an etcd grpc-proxy sidecar!"
  end

  run! if app_file == $0
end