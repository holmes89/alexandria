module Main exposing (main)

import Browser
import Page.ListBooks as ListBooks


main : Program () ListBooks.Model ListBooks.Msg
main =
    Browser.element
        { init = ListBooks.init
        , view = ListBooks.view
        , update = ListBooks.update
        , subscriptions = \_ -> Sub.none
        }
