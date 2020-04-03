module Route exposing (Route(..), parseUrl)

import Book exposing (BookID)
import Url exposing (Url)
import Url.Parser exposing (..)


type Route
    = NotFound
    | Login
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
        [ map Login top
        , map Login (s "login")
        , map Books (s "books")
        , map Book (s "books" </> string)
        ]
