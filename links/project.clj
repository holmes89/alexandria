(defproject links "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :license {:name "EPL-2.0 OR GPL-2.0-or-later WITH Classpath-exception-2.0"
            :url "https://www.eclipse.org/legal/epl-2.0/"}
  :dependencies [[org.clojure/clojure "1.10.1"]
                 [toucan "1.15.0"]
                 [prismatic/schema "1.1.12"]
                 [metosin/compojure-api "2.0.0-alpha31"]
                 [ring/ring-jetty-adapter "1.8.0"]
                 [org.postgresql/postgresql "42.2.11"]
                 [environ "1.1.0"]
                 [buddy "2.0.0"]
                 [migratus "1.2.8"]]
  :main ^:skip-aot links.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
