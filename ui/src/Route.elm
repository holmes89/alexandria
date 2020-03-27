module Route exposing (Route(..), parseUrl)

import Url exposing (Url)
import Url.Parser exposing (..)


type Route
    = NotFound
    | Books


parseUrl : Url -> Route
parseUrl url =
    case parse matchRoute url of
        Just route ->
            route

        Nothing ->
            NotFound


matchRoute : Parser (Route -> a) a
matchRoute =
    oneOf
        [ map Books top
        , map Books (s "books")
        ]
