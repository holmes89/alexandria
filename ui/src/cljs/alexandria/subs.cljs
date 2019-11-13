(ns alexandria.subs
  (:require
   [re-frame.core :as re-frame]))

(re-frame/reg-sub
 ::name
 (fn [db]
   (:name db)))

(re-frame/reg-sub
    ::active-panel
  (fn [db _]
    (:active-panel db)))

(re-frame/reg-sub
    ::docs
  (fn [db _]
    (:document-data db)))

(re-frame/reg-sub
    ::active-doc
  (fn [db _]
    (:active-doc db)))
