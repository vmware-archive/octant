class YouTube < Liquid::Tag
  Syntax = /^\s*([^\s]+)\W"([^"]*)"(\s+(\d+)\s+(\d+)\s*)?/
  def initialize(tagName, markup, tokens)
    super

    if markup =~ Syntax then
      @id = $1
      @title = $2
      if $3.nil? then
          @width = 560
          @height = 420
      else
          @width = $3.to_i
          @height = $4.to_i
      end
    else
      raise "No YouTube ID provided in the \"youtube\" tag"
    end
  end

  def render(context)
    # "<iframe width=\"#{@width}\" height=\"#{@height}\" src=\"https://www.youtube.com/embed/#{@id}\" frameborder=\"0\"allowfullscreen></iframe>"
    "<iframe title=\"#{@title}\" width=\"#{@width}\" height=\"#{@height}\" src=\"https://www.youtube.com/embed/#{@id}?color=white&theme=light\"></iframe>"
  end

  Liquid::Template.register_tag "youtube", self
end
