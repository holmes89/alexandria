module Route exposing (Route(..), parseUrl)

import Book exposing (BookID)
import Url exposing (Url)
import Url.Parser exposing (..)


type Route
    = NotFound
    | Books
    | Book BookID


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
        , map Book (s "books" </> string)
        ]
