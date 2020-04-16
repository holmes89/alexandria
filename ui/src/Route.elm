module Route exposing (Route(..), parseUrl)

import Book exposing (BookID)
import Paper exposing (PaperID)
import Url exposing (Url)
import Url.Parser exposing (..)


type Route
    = NotFound
    | Home
    | Login
    | Books
    | Book BookID
    | Papers
    | Paper PaperID
    | Journal
    | Links


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
        [ map Home top
        , map Login (s "login")
        , map Books (s "books")
        , map Book (s "books" </> string)
        , map Papers (s "papers")
        , map Paper (s "papers" </> string)
        , map Journal (s "journal")
        , map Links (s "links")
        ]
